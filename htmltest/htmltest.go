package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"path"
	"sync"
	"time"
)

type HtmlTest struct {
	opts       Options
	httpClient *http.Client
	documents  []htmldoc.Document
	issueStore issues.IssueStore
}

type HtmlTester interface {
	Test()

	setOptions(map[string]interface{})
	testDocuments()
	parseNode()

	checkLink(document *htmldoc.Document, node *html.Node)
	checkImg(document *htmldoc.Document, node *html.Node)
	checkScript(document *htmldoc.Document, node *html.Node)

	checkExternal(node *html.Node)
	checkInternal(node *html.Node)
	checkFile(node *html.Node, fPath string)
	checkMailto(node *html.Node)
	checkTel(node *html.Node)
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
		hT.documents = htmldoc.DocumentsFromDir(hT.opts.DirectoryPath)
	} else {
		log.Fatal("Neither file or directory path provided")
	}

	hT.testDocuments()

	return &hT
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (hT *HtmlTest) testDocuments() {
	if hT.opts.TestFilesConcurrently {
		var wg sync.WaitGroup
		for _, document := range hT.documents {
			wg.Add(1)
			go func(document htmldoc.Document) {
				defer wg.Done()
				document.Parse()
				hT.parseNode(&document, document.HTMLNode)
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
