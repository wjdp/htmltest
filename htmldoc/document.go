package htmldoc

import (
	"golang.org/x/net/html"
	"os"
	"sync"
)

type Document struct {
	FilePath        string // Relative to the shell session
	SitePath        string // Relative to the site root
	Directory       string
	htmlMutex       *sync.Mutex
	htmlNode        *html.Node
	hashMap         map[string]*html.Node
	NodesOfInterest []*html.Node
	State           DocumentState
}

// Used by checks that depend on the document being parsed
type DocumentState struct {
	FaviconPresent bool
}

func (doc *Document) Init() {
	// Setup the document, doesn't mesh nice with the NewXYZ() convention but
	// many optional parameters for Document and no parameter overloading in Go
	doc.htmlMutex = &sync.Mutex{}
	doc.NodesOfInterest = make([]*html.Node, 0)
	doc.hashMap = make(map[string]*html.Node)
}

func (doc *Document) Parse() {
	// Parse the document
	// Either called when the document is tested or when another document needs
	// data from this one.
	doc.htmlMutex.Lock() // MUTEX
	if doc.htmlNode != nil {
		doc.htmlMutex.Unlock() // MUTEX
		return
	}
	// Open, parse, and close document
	f, err := os.Open(doc.FilePath)
	checkErr(err)
	defer f.Close()

	htmlNode, err := html.Parse(f)
	checkErr(err)

	doc.htmlNode = htmlNode
	doc.parseNode(htmlNode)
	doc.htmlMutex.Unlock() // MUTEX
}

func (doc *Document) parseNode(n *html.Node) {
	if n.Type == html.ElementNode {
		// If present save fragment identifier to the hashMap
		nodeId := GetId(n.Attr)
		if nodeId != "" {
			doc.hashMap[nodeId] = n
		}
		// Identify and store tags of interest
		switch n.Data {
		case "a", "link", "img", "script":
			doc.NodesOfInterest = append(doc.NodesOfInterest, n)
		case "pre", "code":
			return // Everything within these elements is not to be interpreted
		}
	}
	// Iterate over children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		doc.parseNode(c)
	}
}

func (doc *Document) IsHashValid(hash string) bool {
	doc.Parse() // Ensure doc has been parsed
	_, ok := doc.hashMap[hash]
	return ok
}
