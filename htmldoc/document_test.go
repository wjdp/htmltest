package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestDocumentParse(t *testing.T) {
	// parse a document and check we have valid nodes
	doc := Document{
		FilePath: "fixtures/documents/index.html",
	}
	doc.Init()
	doc.Parse()
	nodeElem := doc.htmlNode.FirstChild.FirstChild.NextSibling.FirstChild
	assert.Equals(t, "document first body node", nodeElem.Data, "h1")
}

func TestDocumentNodesOfInterest(t *testing.T) {
	doc := Document{
		FilePath: "fixtures/documents/nodes.htm",
	}
	doc.Init()
	doc.Parse()
	assert.Equals(t, "nodes of interest", len(doc.NodesOfInterest), 4)
}

func TestDocumentIsHashValid(t *testing.T) {
	// parse a document and check we have valid nodes
	doc := Document{
		FilePath: "fixtures/documents/index.html",
	}
	doc.Init()
	doc.Parse()

	assert.IsTrue(t, "#xyz present", doc.IsHashValid("xyz"))
	assert.IsTrue(t, "#prq present", doc.IsHashValid("prq"))
	assert.IsFalse(t, "#abc present", doc.IsHashValid("abc"))
}
