package htmldoc

import (
	"golang.org/x/net/html"
	"os"
)

type Document struct {
	FilePath  string // Relative to the shell session
	SitePath  string // Relative to the site root
	Directory string
	HTMLNode  *html.Node
	State     DocumentState
}

// Used by checks that depend on the document being parsed
type DocumentState struct {
	FaviconPresent bool
}

func (doc *Document) Parse() {
	// Open, parse, and close document
	f, err := os.Open(doc.FilePath)
	checkErr(err)
	defer f.Close()

	htmlNode, err := html.Parse(f)
	checkErr(err)

	doc.HTMLNode = htmlNode
}
