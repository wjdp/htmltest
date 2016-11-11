package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/refcache"
	"golang.org/x/net/html"
	"net/http"
	"path"
	"sync"
	"time"
)

type HtmlTest struct {
	opts        Options
	httpClient  *http.Client
	httpChannel chan bool
	documents   []htmldoc.Document
	issueStore  issues.IssueStore
	refCache    *refcache.RefCache
}

func Test(optsUser map[string]interface{}) *HtmlTest {
	hT := HtmlTest{}

	// Merge user options with defaults and set hT.opts
	hT.setOptions(optsUser)

	// Create issue store and set LogLevel
	hT.issueStore = issues.NewIssueStore(hT.opts.LogLevel)

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

	if hT.opts.FilePath != "" {
		// Single document mode
		doc := htmldoc.Document{
			FilePath: path.Join(hT.opts.DirectoryPath, hT.opts.FilePath),
			SitePath: hT.opts.FilePath,
		}
		hT.documents = []htmldoc.Document{doc}
	} else if hT.opts.DirectoryPath != "" {
		// Directory mode
		hT.documents = htmldoc.DocumentsFromDir(
			hT.opts.DirectoryPath, hT.opts.IgnoreDirs)
	} else {
		panic("Neither file or directory path provided")
	}

	hT.testDocuments()

	if hT.opts.EnableCache {
		hT.refCache.WriteStore(cachePath)
	}
	if hT.opts.EnableLog {
		hT.issueStore.WriteLog(path.Join(hT.opts.ProgDir, hT.opts.LogFile))
	}

	return &hT
}

func (hT *HtmlTest) testDocuments() {
	if hT.opts.TestFilesConcurrently {
		hT.issueStore.AddIssue(issues.Issue{
			Level:   issues.WARNING,
			Message: "running in concurrent mode, this is experimental",
		})
		var wg sync.WaitGroup
		// Make buffered channel to act as concurrency limiter
		var concChannel = make(chan bool, hT.opts.DocumentConcurrencyLimit)
		for _, document := range hT.documents {
			wg.Add(1)
			concChannel <- true // Add to concurrency limiter
			go func(document htmldoc.Document) {
				defer wg.Done()
				hT.testDocument(&document)
				<-concChannel // Bump off concurrency limiter
			}(document)
		}
		wg.Wait()
	} else {
		for _, document := range hT.documents {
			hT.testDocument(&document)
		}
	}
}

func (hT *HtmlTest) testDocument(document *htmldoc.Document) {
	document.Parse()
	hT.parseNode(document, document.HTMLNode)
	hT.postChecks(document)
}

func (hT *HtmlTest) parseNode(document *htmldoc.Document, n *html.Node) {
	if n.Type == html.ElementNode {
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
		case "pre":
			return // Everything within a pre is not to be interpreted
		case "code":
			return // Everything within a code is not to be interpreted
		}
	}
	// Iterate over children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		hT.parseNode(document, c)
	}
}

func (hT *HtmlTest) postChecks(document *htmldoc.Document) {
	// Checks to run after document has been parsed
	if hT.opts.CheckFavicon && !document.State.FaviconPresent {
		hT.issueStore.AddIssue(issues.Issue{
			Level:   issues.ERROR,
			Message: "favicon missing",
		})
	}
}

func (hT *HtmlTest) CountErrors() int {
	return hT.issueStore.Count(issues.ERROR)
}
