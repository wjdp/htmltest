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
	doc.Parse()
	nodeElem := doc.HTMLNode.FirstChild.FirstChild.NextSibling.FirstChild
	assert.Equals(t, "document first body node", nodeElem.Data, "h1")
}
