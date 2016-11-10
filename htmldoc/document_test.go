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

func TestDocumentsFromDir(t *testing.T) {
	// it creates Document struts from an os directory
	docs := DocumentsFromDir("fixtures/documents", []interface{}{})
	// Fixtures dir has seven documents in various folders
	assert.Equals(t, "document count", len(docs), 7)
}
