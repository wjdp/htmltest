// Package htmltest : Main package, provides the HTMLTest struct and
// associated checks.
package htmltest

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/output"
	"github.com/wjdp/htmltest/refcache"
	"gopkg.in/seborama/govcr.v4"
)

// Base path for VCR cassettes, relative to this package
const vcrCassetteBasePath string = "fixtures/vcr"

// HTMLTest struct, A html testing session, user options are passed in and
// tests are run.
type HTMLTest struct {
	opts          Options
	httpClient    *http.Client
	httpChannel   chan bool
	documentStore htmldoc.DocumentStore
	issueStore    issues.IssueStore
	refCache      *refcache.RefCache
}

func setRedirectLimitCheck(hT HTMLTest) func(req *http.Request, via []*http.Request) error {
	redirectLimit := hT.opts.RedirectLimit

	// Nothing set or invalid, use defaults from net/http
	if 0 > redirectLimit {
		return nil
	}

	return func(req *http.Request, via []*http.Request) error {
		if redirectLimit < len(via) {
			originalURL := via[0].URL.String()
			hT.issueStore.AddIssue(issues.Issue{
				Level:   issues.LevelError,
				Message: "too many redirects: " + originalURL,
			})
			return errors.New("too many redirects: " + originalURL)
		}
		return nil
	}
}

// Test : Given user options run htmltest and return a pointer to the test
// object.
func Test(optsUser map[string]interface{}) (*HTMLTest, error) {
	hT := HTMLTest{}

	// If FilePath set, modify FileExtension
	if optsUser["FilePath"] != nil {
		optsUser["FileExtension"] = path.Ext(optsUser["FilePath"].(string))
	}

	// Merge user options with defaults and set hT.opts
	hT.setOptions(optsUser)

	// Create issue store and set LogLevel and printImmediately if sort is seq
	hT.issueStore = issues.NewIssueStore(hT.opts.LogLevel,
		(hT.opts.LogSort == "seq"))

	transport := &http.Transport{
		// Disable HTTP/2, this is required due to a number of edge cases where http negotiates H2, but something goes
		// wrong when actually using it. Downgrading to H1 when this issue is hit is not yet supported so we use the
		// following to disable H2 support:
		// > Programs that must disable HTTP/2 can do so by setting Transport.TLSNextProto ... to a non-nil, empty map.
		// See issue #49
		TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: hT.opts.IgnoreSSLVerify},
	}
	hT.httpClient = &http.Client{
		// Durations are in nanoseconds
		Transport:     transport,
		Timeout:       time.Duration(hT.opts.ExternalTimeout) * time.Second,
		CheckRedirect: setRedirectLimitCheck(hT),
	}

	// If enabled (unit tests only) patch in govcr to the httpClient
	var vcr *govcr.VCRControlPanel
	if hT.opts.VCREnable {
		// Strip fixtures/ from the start of the path. This will break if the path doesn't start with "fixtures/"
		cassettePath := strings.Split(hT.opts.DirectoryPath, "fixtures/")[1]
		// Build VCR
		vcr = govcr.NewVCR(hT.opts.FilePath,
			&govcr.VCRConfig{
				Client:       hT.httpClient,
				CassettePath: path.Join(vcrCassetteBasePath, cassettePath),
			})

		// Inject VCR's http.Client wrapper
		hT.httpClient = vcr.Client
	}

	// Make buffered channel to act as concurrency limiter
	hT.httpChannel = make(chan bool, hT.opts.HTTPConcurrencyLimit)

	// Setup refCache
	cachePath := ""
	if hT.opts.EnableCache {
		cachePath = path.Join(hT.opts.OutputDir, hT.opts.OutputCacheFile)
	}
	hT.refCache = refcache.NewRefCache(cachePath, hT.opts.CacheExpires)

	if hT.opts.NoRun {
		return &hT, nil
	}

	// Either of these options are required to run
	if hT.opts.DirectoryPath == "" && hT.opts.FilePath == "" {
		err := errors.New("Neither FilePath nor DirectoryPath provided")
		return &hT, err
	}

	// Check the provided DirectoryPath exists
	f, err := os.Open(hT.opts.DirectoryPath)
	if os.IsNotExist(err) {
		err := errors.New(fmt.Sprint(
			"Cannot access '" + hT.opts.DirectoryPath + "', no such directory."))
		return &hT, err
	}
	// Get FileInfo, (scan for details)
	fi, err := f.Stat()
	output.CheckErrorPanic(err)
	// Check if DirectoryPath directory
	if !fi.IsDir() {
		err := errors.New(fmt.Sprint(
			"DirectoryPath '" + hT.opts.DirectoryPath + "' is a file, not a directory."))
		return &hT, err
	}

	// Init our document store
	hT.documentStore = htmldoc.NewDocumentStore()
	// Setup document store
	hT.documentStore.BasePath = hT.opts.DirectoryPath
	hT.documentStore.DocumentExtension = hT.opts.FileExtension
	hT.documentStore.DirectoryIndex = hT.opts.DirectoryIndex
	hT.documentStore.IgnorePatterns = hT.opts.IgnoreDirs
	hT.documentStore.IgnoreTagAttribute = hT.opts.IgnoreTagAttribute

	if hT.opts.BaseURL != "" {
		baseURL, err := url.Parse(hT.opts.BaseURL)
		if err != nil {
			err := fmt.Errorf("Could not parse BaseURL '%s': %w", hT.opts.BaseURL, err)
			return &hT, err
		}

		hT.documentStore.BaseURL = baseURL
	}

	// Discover documents
	hT.documentStore.Discover()

	if hT.opts.FilePath != "" {
		// Single document mode
		doc, ok := hT.documentStore.ResolvePath(hT.opts.FilePath)
		if !ok {
			err := errors.New(fmt.Sprint(
				"Could not find FilePath '", hT.opts.FilePath, "' in '", hT.opts.DirectoryPath, "'"))
			return &hT, err
		}
		hT.testDocument(doc)
	} else if hT.opts.DirectoryPath != "" {
		// Test documents
		hT.testDocuments()
	}

	if hT.opts.EnableCache {
		hT.refCache.WriteStore(cachePath)
	}
	if hT.opts.EnableLog {
		hT.issueStore.WriteLog(path.Join(hT.opts.OutputDir,
			hT.opts.OutputLogFile))
	}

	// This is useful for debugging the VCR, but rather noisy otherwise
	//if hT.opts.VCREnable {
	//	fmt.Printf("%+v\n", vcr.Stats())
	//}

	return &hT, nil
}

func (hT *HTMLTest) testDocuments() {
	if hT.opts.TestFilesConcurrently {
		hT.issueStore.AddIssue(issues.Issue{
			Level:   issues.LevelWarning,
			Message: "running in concurrent mode, this is experimental",
		})
		var wg sync.WaitGroup
		// Make buffered channel to act as concurrency limiter
		var concChannel = make(chan bool, hT.opts.DocumentConcurrencyLimit)
		for _, document := range hT.documentStore.Documents {
			wg.Add(1)
			concChannel <- true // Add to concurrency limiter
			go func(document *htmldoc.Document) {
				defer wg.Done()
				hT.testDocument(document)
				<-concChannel // Bump off concurrency limiter
			}(document)
		}
		wg.Wait()
	} else {
		for _, document := range hT.documentStore.Documents {
			hT.testDocument(document)
		}
	}
}

func (hT *HTMLTest) testDocument(document *htmldoc.Document) {
	if document.IgnoreTest {
		hT.issueStore.AddIssue(issues.Issue{
			Level:   issues.LevelDebug,
			Message: "ignored " + document.SitePath,
		})
		return
	}

	hT.issueStore.AddIssue(issues.Issue{
		Level:   issues.LevelDebug,
		Message: "testDocument on " + document.SitePath,
	})

	document.Parse()

	if hT.opts.CheckDoctype {
		hT.checkDoctype(document)
	}

	for _, n := range document.NodesOfInterest {
		switch n.Data {
		case "a":
			if hT.opts.CheckAnchors {
				hT.checkLink(document, n)
			}
		case "link":
			if hT.opts.CheckLinks {
				hT.checkLink(document, n)
			}
		case "img":
			if hT.opts.CheckImages {
				hT.checkImg(document, n)
			}
		case "script":
			if hT.opts.CheckScripts {
				hT.checkScript(document, n)
			}
		case "meta":
			if hT.opts.CheckMeta {
				hT.checkMeta(document, n)
			}
		case "area":
			if hT.opts.CheckGeneric {
				hT.checkGeneric(document, n, "href")
			}
		case "blockquote", "del", "ins", "q":
			if hT.opts.CheckGeneric {
				hT.checkGeneric(document, n, "cite")
			}
		case "iframe", "input", "audio", "embed", "source", "track":
			if hT.opts.CheckGeneric {
				hT.checkGeneric(document, n, "src")
			}
		case "video":
			if hT.opts.CheckGeneric {
				hT.checkGeneric(document, n, "src")
				hT.checkGeneric(document, n, "poster")
			}
		case "object":
			if hT.opts.CheckGeneric {
				hT.checkGeneric(document, n, "data")
			}
		}
	}
	hT.postChecks(document)

	// If sorting by document output issues now
	if hT.opts.LogSort == "document" {
		hT.issueStore.PrintDocumentIssues(document)
	}
}

func (hT *HTMLTest) postChecks(document *htmldoc.Document) {
	// Checks to run after document has been parsed
	if hT.opts.CheckFavicon && !document.State.FaviconPresent {
		hT.issueStore.AddIssue(issues.Issue{
			Level:   issues.LevelError,
			Message: "favicon missing",
		})
	}
}

// CountErrors : Return number of error level issues
func (hT *HTMLTest) CountErrors() int {
	return hT.issueStore.Count(issues.LevelError)
}

// CountDocuments : Return number of documents in hT document store
func (hT *HTMLTest) CountDocuments() int {
	return len(hT.documentStore.Documents)
}
