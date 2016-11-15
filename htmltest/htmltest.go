// Package htmltest : Main package, provides the HTMLTest struct and
// associated checks.
package htmltest

import (
	"fmt"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/refcache"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

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

// Test : Given user options run htmltest and return a pointer to the test
// object.
func Test(optsUser map[string]interface{}) *HTMLTest {
	hT := HTMLTest{}

	// Merge user options with defaults and set hT.opts
	hT.setOptions(optsUser)

	// Create issue store and set LogLevel and printImmediately if sort is "seq"
	hT.issueStore = issues.NewIssueStore(hT.opts.LogLevel, (hT.opts.LogSort == "seq"))

	transport := &http.Transport{
		TLSNextProto: nil, // Disable HTTP/2, "write on closed buffer" errors
	}
	hT.httpClient = &http.Client{
		// Durations are in nanoseconds
		Transport: transport,
		Timeout:   time.Duration(hT.opts.ExternalTimeout * 1000000000),
	}

	// Make buffered channel to act as concurrency limiter
	hT.httpChannel = make(chan bool, hT.opts.HTTPConcurrencyLimit)

	// Setup refcache
	cachePath := ""
	if hT.opts.EnableCache {
		cachePath = path.Join(hT.opts.ProgDir, hT.opts.CacheFile)
	}
	hT.refCache = refcache.NewRefCache(cachePath, hT.opts.CacheExpires)

	if hT.opts.NoRun {
		return &hT
	}

	// Init our document store
	hT.documentStore = htmldoc.NewDocumentStore()
	// Setup document store
	hT.documentStore.BasePath = hT.opts.DirectoryPath
	hT.documentStore.DocumentExtension = ".html" // TODO add option
	hT.documentStore.DirectoryIndex = hT.opts.DirectoryIndex
	hT.documentStore.IgnorePatterns = hT.opts.IgnoreDirs
	// Discover documents
	hT.documentStore.Discover()

	if hT.opts.FilePath != "" {
		// Single document mode
		doc, ok := hT.documentStore.ResolvePath(hT.opts.FilePath)
		if !ok {
			fmt.Println("Could not find document", hT.opts.FilePath, "in", hT.opts.DirectoryPath)
			os.Exit(1)
		}
		hT.testDocument(doc)
	} else if hT.opts.DirectoryPath != "" {
		// Test documents
		hT.testDocuments()
	} else {
		panic("Neither file or directory path provided")
	}

	if hT.opts.EnableCache {
		hT.refCache.WriteStore(cachePath)
	}
	if hT.opts.EnableLog {
		hT.issueStore.WriteLog(path.Join(hT.opts.ProgDir, hT.opts.LogFile))
	}

	return &hT
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
	document.Parse()
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
