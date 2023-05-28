package htmldoc

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/theunrepentantgeek/htmltest/output"
	"golang.org/x/net/html"
)

// Document struct, representation of a document within the tested site
type Document struct {
	FilePath           string                // Relative to the shell session
	SitePath           string                // Relative to the site root
	BasePath           string                // Base for relative links
	IgnoreTest         bool                  // Ignore this Document for testing.
	htmlMutex          *sync.Mutex           // Controls access to htmlNode
	htmlNode           *html.Node            // Parsed output
	hashMap            map[string]*html.Node // Map of valid id/names of nodes
	NodesOfInterest    []*html.Node          // Slice of nodes to run checks on
	State              DocumentState         // Link to a DocumentState struct
	DoctypeNode        *html.Node            // Pointer to doctype node if exists
	ignoreTagAttribute string                // Attribute to ignore element and children if found on element
}

// DocumentState struct, used by checks that depend on the document being
// parsed.
type DocumentState struct {
	FaviconPresent bool // Have we found a favicon in the document?
}

// Init : Initialise the Document struct doesn't mesh nice with the NewXYZ()
// convention but many optional parameters for Document and no parameter
// overloading in Go
func (doc *Document) Init() {
	// Setup the document,
	doc.htmlMutex = &sync.Mutex{}
	doc.NodesOfInterest = make([]*html.Node, 0)
	doc.hashMap = make(map[string]*html.Node)
}

// Parse : Ask Document to parse its HTML file. Returns quickly if this has
// already been done. Thread safe. Either called when the document is tested
// or when another document needs data from this one.
func (doc *Document) Parse() {

	// Only one routine may parse the doc
	doc.htmlMutex.Lock()
	defer doc.htmlMutex.Unlock()

	// If document has already been parsed, return early.
	if doc.htmlNode != nil {
		return
	}

	// Open, parse, and close document
	f, err := os.Open(doc.FilePath)
	output.CheckErrorPanic(err)
	defer f.Close()

	htmlNode, err := html.Parse(f)
	output.CheckErrorGeneric(err)

	doc.htmlNode = htmlNode
	doc.parseNode(htmlNode)
}

// Internal recursive function that delves into the node tree and captures
// nodes of interest and node id/names.
func (doc *Document) parseNode(n *html.Node) {
	// Ignore this tree if data-proofer-ignore set
	if doc.ignoreTagAttribute != "" &&
		(AttrPresent(n.Attr, doc.ignoreTagAttribute) || ClassPresent(n.Attr, doc.ignoreTagAttribute)) {
		return
	}

	switch n.Type {
	case html.DoctypeNode:
		doc.DoctypeNode = n
	case html.ElementNode:
		// If present save fragment identifier to the hashMap
		nodeID := GetID(n.Attr)
		if nodeID != "" {
			doc.hashMap[nodeID] = n
		}
		// Identify and store tags of interest
		switch n.Data {
		case "a", "area", "audio", "blockquote", "del", "embed", "iframe", "img",
			"input", "ins", "link", "meta", "object", "q", "script", "source",
			"track", "video":
			// Nodes of interest
			doc.NodesOfInterest = append(doc.NodesOfInterest, n)
		case "base":
			// Set BasePath from <base> tag
			doc.BasePath = path.Join(doc.BasePath, GetAttr(n.Attr, "href"))
		}
	case html.ErrorNode:
		fmt.Printf("%+v\n", n)
		fmt.Println("Oops, in parsing your HTML we fell over.\n",
			"Please let the developer know about this.\n",
			"https://github.com/theunrepentantgeek/htmltest/issues/new")
	}

	// Iterate over children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		doc.parseNode(c)
	}
}

// IsHashValid : Is a hash/fragment present in this Document.
func (doc *Document) IsHashValid(hash string) bool {
	doc.Parse() // Ensure doc has been parsed
	_, ok := doc.hashMap[hash]
	return ok
}
