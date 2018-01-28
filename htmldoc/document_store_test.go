package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestDocumentStoreDiscover(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html" // Ignores .htm
	dS.DirectoryIndex = "index.html"
	dS.Discover()
	// Fixtures dir has eight documents in various folders
	assert.Equals(t, "document count", len(dS.Documents), 6)
}

func TestDocumentStoreIgnorePatterns(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html" // Ignores .htm
	dS.DirectoryIndex = "index.html"
	dS.IgnorePatterns = []string{"^lib/"}
	dS.Discover()
	// Fixtures dir has seven documents in various folders, (one ignored in lib)
	assert.Equals(t, "document count", len(dS.Documents), 5)
}

func TestDocumentStoreDocumentExists(t *testing.T) {
	// documentstore knows if documents exist or not
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html"
	dS.DirectoryIndex = "index.html"
	dS.Discover()
	_, b1 := dS.DocumentPathMap["index.html"]
	assert.IsTrue(t, "index.html exists", b1)
	_, b2 := dS.DocumentPathMap["dir2/index.html"]
	assert.IsTrue(t, "dir2/index.html exists", b2)
	_, b3 := dS.DocumentPathMap["foo.html"]
	assert.IsFalse(t, "foo.html does not exist", b3)
	_, b4 := dS.DocumentPathMap["dir3/index.html"]
	assert.IsFalse(t, "dir3/index.html does not exist", b4)
}

func TestDocumentStoreDocumentResolve(t *testing.T) {
	// documentstore correctly resolves documents
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html"
	dS.DirectoryIndex = "index.html"
	dS.Discover()
	d0, b0 := dS.ResolvePath("/")
	assert.IsTrue(t, "root document exists", b0)
	assert.Equals(t, "/ resolves to index.html",
		d0.FilePath, "fixtures/documents/index.html")
	d1, b1 := dS.ResolvePath("/contact.html")
	assert.IsTrue(t, "/contact.html exists", b1)
	assert.Equals(t, "/contact.html resolves to correct document",
		d1.FilePath, "fixtures/documents/contact.html")
	d2, b2 := dS.ResolvePath("dir2/index.html")
	assert.IsTrue(t, "dir2/index.html exists", b2)
	assert.Equals(t, "dir2/index.html resolves to correct document",
		d2.FilePath, "fixtures/documents/dir2/index.html")
	d3, b3 := dS.ResolvePath("dir2/")
	assert.IsTrue(t, "dir2/index.html exists", b3)
	assert.Equals(t, "dir2/index.html resolves to correct document",
		d3.FilePath, "fixtures/documents/dir2/index.html")
	d4, b4 := dS.ResolvePath("dir2")
	assert.IsTrue(t, "dir2/index.html exists", b4)
	assert.Equals(t, "dir2/index.html resolves to correct document",
		d4.FilePath, "fixtures/documents/dir2/index.html")
	_, b5 := dS.ResolvePath("does-not-exist")
	assert.IsFalse(t, "does not return doc for invalid path", b5)
}
