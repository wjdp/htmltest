package htmldoc

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/wjdp/htmltest/output"
	"golang.org/x/net/html"
)

// Document struct, representation of a document within the tested site
type Document struct {
	FilePath           string              // Relative to the shell session
	SitePath           string              // Relative to the site root
	BasePath           string              // Base for relative links
	IgnoreTest         bool                // Ignore this Document for testing.
	htmlMutex          *sync.Mutex         // Controls access to htmlNode
	htmlNode           *html.Node          // Parsed output
	parsed             bool                // Have hashMap/DoctypeNode been built at least once?
	treeLoaded         bool                // Is the node tree currently in memory?
	basePathOrig       string              // BasePath as discovered, before any <base> tag was applied
	hashMap            map[string]struct{} // Set of valid id/names of nodes
	collectHashes      bool                // Record id/names? Only needed for hash checks
	NodesOfInterest    []*html.Node        // Slice of nodes to run checks on
	State              DocumentState       // Link to a DocumentState struct
	DoctypeNode        *html.Node          // Pointer to doctype node if exists
	ignoreTagAttribute string              // Attribute to ignore element and children if found on element
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
	doc.hashMap = make(map[string]struct{})
	// Default to collecting hashes so a Document used directly behaves as
	// before; DocumentStore overrides this from the CheckInternalHash option.
	doc.collectHashes = true
}

// Parse : Ask Document to parse its HTML file. Returns quickly if this has
// already been done. Thread safe. Either called when the document is tested
// or when another document needs data from this one.
func (doc *Document) Parse() {
	// Only one routine may parse the doc
	doc.htmlMutex.Lock()
	defer doc.htmlMutex.Unlock()

	doc.parseLocked()
}

// parseLocked : Parse body, caller must hold htmlMutex.
func (doc *Document) parseLocked() {
	// If the tree is already in memory, return early.
	if doc.treeLoaded {
		return
	}

	// parseNode derives BasePath from any <base> tag by joining onto the
	// current value, so the discovered value is restored first; otherwise a
	// re-parse after Release would apply the <base> href a second time.
	if doc.parsed {
		doc.BasePath = doc.basePathOrig
	} else {
		doc.basePathOrig = doc.BasePath
	}

	// Open, parse, and close document
	f, err := os.Open(doc.FilePath)
	output.CheckErrorPanic(err)
	defer f.Close()

	htmlNode, err := html.Parse(f)
	output.CheckErrorGeneric(err)

	doc.htmlNode = htmlNode
	// Reset so a re-parse after Release cannot double up on the previous run.
	doc.NodesOfInterest = make([]*html.Node, 0)
	doc.parseNode(htmlNode)
	doc.treeLoaded = true
	doc.parsed = true
}

// Release : Discard the parsed node tree once this document's own checks are
// done. Every html.Node holds Parent/child/sibling pointers, so retaining a
// single node keeps that document's entire tree reachable. Without this the
// whole site's DOM accumulates for the lifetime of the run.
//
// hashMap is deliberately kept: it holds only strings, and other documents may
// still need it to validate fragment links against this one.
func (doc *Document) Release() {
	doc.htmlMutex.Lock()
	defer doc.htmlMutex.Unlock()

	doc.releaseLocked()
}

// releaseLocked : Release body, caller must hold htmlMutex.
func (doc *Document) releaseLocked() {
	doc.htmlNode = nil
	doc.DoctypeNode = nil
	doc.NodesOfInterest = nil
	doc.treeLoaded = false
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
		// If present save fragment identifier to the hashMap. Only the
		// presence of the id is ever read, so the node itself is not stored:
		// every html.Node links to its parent and siblings, so keeping one
		// keeps the document's whole tree reachable.
		if doc.collectHashes {
			nodeID := GetID(n.Attr)
			if nodeID != "" {
				doc.hashMap[nodeID] = struct{}{}
			}
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
			"https://github.com/wjdp/htmltest/issues/new")
	}

	// Iterate over children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		doc.parseNode(c)
	}
}

// IsHashValid : Is a hash/fragment present in this Document.
func (doc *Document) IsHashValid(hash string) bool {
	doc.htmlMutex.Lock()
	defer doc.htmlMutex.Unlock()

	// A document is commonly parsed early because another document links to a
	// fragment in it. Load it just far enough to learn its fragments and drop
	// the tree again, otherwise a site that cross-links fragments would hold a
	// large part of its DOM at once.
	//
	// Load and release happen under a single hold of the mutex: nobody else
	// can observe the tree in between, so this cannot release a tree that
	// another routine is testing. !parsed implies !treeLoaded, so the only
	// tree released here is the one loaded here.
	if !doc.parsed {
		doc.parseLocked()
		doc.releaseLocked()
	}

	_, ok := doc.hashMap[hash]
	return ok
}
