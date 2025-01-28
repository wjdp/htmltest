package htmldoc

import (
	"testing"

	"github.com/daviddengcn/go-assert"
)

// Expected number of .html files under "fixtures/documents"
const ExpectedHtmlDocumentCount = 6

func TestDocumentStoreDiscover(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html" // Ignores .htm
	dS.DirectoryIndex = "index.html"
	dS.Discover()
	assert.Equals(t, "document count", dS.DocumentCount(), ExpectedHtmlDocumentCount)
	assert.Equals(t, "ignored document count", dS.IgnoredDocCount(), 0)

	for _, document := range dS.Documents {
		assert.IsFalse(t, document.SitePath+" is not ignored", document.IgnoreTest)
	}
}

func TestDocumentStoreIgnorePatterns(t *testing.T) {
	// documentstore can scan an os directory
	dS := NewDocumentStore()
	dS.BasePath = "fixtures/documents"
	dS.DocumentExtension = ".html" // Ignores .htm
	dS.DirectoryIndex = "index.html"
	dS.IgnorePatterns = []interface{}{"^lib/"}
	dS.Discover()
	// IgnorePatterns does not affect stored document count
	assert.Equals(t, "document count", dS.DocumentCount(), ExpectedHtmlDocumentCount)
	assert.Equals(t, "ignored document count", dS.IgnoredDocCount(), 1)

	ignoredFile := "lib/unwanted-file.html"
	f, exists := dS.DocumentPathMap[ignoredFile]
	assert.IsTrue(t, ignoredFile+" exists", exists)
	assert.IsTrue(t, ignoredFile+" is flagged as ignored", f.IgnoreTest)

	for _, document := range dS.Documents {
		if document.SitePath != ignoredFile {
			assert.IsFalse(t, document.FilePath+" is not ignored", document.IgnoreTest)
		}
	}
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
