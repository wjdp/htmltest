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
	hT.httpChannel = make(chan bool, 1)

	// Setup refcache
	hT.refCache = refcache.NewRefCache(
		path.Join(hT.opts.ProgDir, hT.opts.CacheFile), hT.opts.CacheExpires)

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

	hT.refCache.WriteStore(path.Join(hT.opts.ProgDir, hT.opts.CacheFile))
	hT.issueStore.WriteLog(path.Join(hT.opts.ProgDir, hT.opts.LogFile))

	return &hT
}

func (hT *HtmlTest) testDocuments() {
	if hT.opts.TestFilesConcurrently {
		var wg sync.WaitGroup
		// Make buffered channel to act as concurrency limiter
		var concChannel = make(chan bool, hT.opts.DocumentConcurrencyLimit)
		for _, document := range hT.documents {
			wg.Add(1)
			concChannel <- true // Add to concurrency limiter
			go func(document htmldoc.Document) {
				defer wg.Done()
				document.Parse()
				hT.parseNode(&document, document.HTMLNode)
				<-concChannel // Bump off concurrency limiter
			}(document)
		}
		wg.Wait()
	} else {
		for _, document := range hT.documents {
			document.Parse()
			hT.parseNode(&document, document.HTMLNode)
		}
	}
}

func (hT *HtmlTest) parseNode(document *htmldoc.Document, n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			hT.checkLink(document, n)
		case "img":
			hT.checkImg(document, n)
		case "link":
			hT.checkLink(document, n)
		case "script":
			hT.checkScript(document, n)
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

func (hT *HtmlTest) CountErrors() int {
	return hT.issueStore.Count(issues.ERROR)
}
