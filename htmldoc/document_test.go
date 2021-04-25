package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"sync"
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

func TestDocumentParseOnce(t *testing.T) {
	// Document.Parse should only parse once, subsequent calls should return quickly
	doc := Document{
		FilePath: "fixtures/documents/index.html",
	}
	doc.Init()
	doc.Parse()
	// Store copy of htmlNode
	hN := doc.htmlNode
	doc.Parse()
	// and assert it's the same one
	assert.Equals(t, "htmlNode", doc.htmlNode, hN)
}

func TestDocumentParseOnceConcurrent(t *testing.T) {
	// Document.Parse should be thread safe
	doc := Document{
		FilePath: "fixtures/documents/index.html",
	}
	doc.Init()
	// Parse many times
	wg := sync.WaitGroup{}
	for i := 0; i < 320; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc.Parse()
		}()
	}
	// Wait until all jobs done
	wg.Wait()
	// Assert we have something sensible by the end of this
	nodeElem := doc.htmlNode.FirstChild.FirstChild.NextSibling.FirstChild
	assert.Equals(t, "document first body node", nodeElem.Data, "h1")
}

func TestDocumentNodesOfInterest(t *testing.T) {
	doc := Document{
		FilePath: "fixtures/documents/nodes.htm",
	}
	doc.Init()
	doc.Parse()
	assert.Equals(t, "nodes of interest", len(doc.NodesOfInterest), 12)
}

func TestDocumentBasePathDefault(t *testing.T) {
	doc := Document{
		FilePath: "fixtures/documents/index.html",
	}
	doc.Init()
	doc.Parse()
	assert.Equals(t, "BasePath", doc.BasePath, "")
}

func TestDocumentBasePathFromTag(t *testing.T) {
	doc := Document{
		FilePath: "fixtures/documents/dir2/base_tag.htm",
	}
	doc.Init()
	doc.Parse()
	assert.Equals(t, "BasePath", doc.BasePath, "/dir2")
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
