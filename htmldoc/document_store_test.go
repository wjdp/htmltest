package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestDocumentStoreDiscover(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = "html"
	dS.Discover()
	// Fixtures dir has eight documents in various folders
	assert.Equals(t, "document count", len(dS.Documents), 8)
}

func TestDocumentStoreIgnorePatterns(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = "html"
	dS.IgnorePatterns = []interface{}{"^lib/"}
	dS.Discover()
	// Fixtures dir has seven documents in various folders, (one ignored in lib)
	assert.Equals(t, "document count", len(dS.Documents), 7)
}

func TestDocumentStoreDocumentExists(t *testing.T) {
	// documentstore knows if documents exist or not
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = "html"
	dS.Discover()
	assert.IsTrue(t, "index.html exists",
		dS.DocumentExists("index.html"))
	assert.IsTrue(t, "dir2/index.html exists",
		dS.DocumentExists("dir2/index.html"))
	assert.IsFalse(t, "foo.html does not exist",
		dS.DocumentExists("foo.html"))
	assert.IsFalse(t, "dir3/index.html does not exist",
		dS.DocumentExists("dir3/index.html"))
}
