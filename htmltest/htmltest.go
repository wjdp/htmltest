package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

var httpClient *http.Client

func setup() {
	issues.LogLevel = Opts.LogLevel
	transport := &http.Transport{
		TLSNextProto: nil, // Disable HTTP/2, "write on closed buffer" errors
	}
	httpClient = &http.Client{
		// Durations are in nanoseconds
		Transport: transport,
		Timeout:   time.Duration(Opts.ExternalTimeout * 1000000000),
	}
}

func Test(optsUser map[string]interface{}) {
	SetOptions(optsUser)
	setup() // Setup objects requiring options
	issues.InitIssueStore()

	if Opts.FilePath != "" {
		doc := htmldoc.Document{
			// Directory: Opts.DirectoryPath,
			Path: Opts.FilePath,
		}
		TestFile(&doc)
	} else if Opts.DirectoryPath != "" {
		TestDirectory(Opts)
	} else {
		log.Fatal("Neither file or directory path provided")
	}
}

func makePath(p string) string {
	return path.Join(Opts.DirectoryPath, p)
}

func TestDirectory(opts Options) {
	log.Printf("htmltest started on %s", Opts.DirectoryPath)

	files := RecurseDirectory("")
	TestFiles(files)
	// issues.OutputIssues()

	log.Printf("%d files checked", len(files))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestFiles(documents []htmldoc.Document) {
	if Opts.TestFilesConcurrently {
		var wg sync.WaitGroup
		for _, document := range documents {
			wg.Add(1)
			go func(document htmldoc.Document) {
				defer wg.Done()
				TestFile(&document)
			}(document)
		}
		wg.Wait()
	} else {
		for _, document := range documents {
			TestFile(&document)
		}
	}
}

func TestFile(document *htmldoc.Document) {
	// log.Println("testFile", document.Path)
	f, err := os.Open(makePath(document.Path))
	checkErr(err)
	defer f.Close()

	document.File = f

	parseHtml(document)
}

func parseHtml(document *htmldoc.Document) {
	doc, err := html.Parse(document.File)
	checkErr(err)
	document.HTMLNode = doc
	parseNode(document, document.HTMLNode)
}

func parseNode(document *htmldoc.Document, n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			CheckLink(document, n)
		case "img":
			CheckImg(document, n)
		case "link":
			CheckLink(document, n)
		case "script":
			CheckScript(document, n)
		case "pre":
			return // Everything within a pre is not to be interpreted
		case "code":
			return // Everything within a code is not to be interpreted
		}
	}
	// Iterate over children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(document, c)
	}
}
