package htmldoc

import "testing"

// A document released after a hash-only parse is parsed again when its own
// checks run. That second parse must reproduce the first, not compound it.
func TestDocumentReparseAfterRelease(t *testing.T) {
	doc := Document{
		FilePath: "fixtures/reparse/base.html",
		SitePath: "base.html",
		BasePath: "dir1",
	}
	doc.Init()

	doc.Parse()
	basePath := doc.BasePath
	nodes := len(doc.NodesOfInterest)
	hashes := len(doc.hashMap)

	// Mimic IsHashValid loading the document then dropping its tree.
	doc.Release()
	if doc.htmlNode != nil || doc.NodesOfInterest != nil {
		t.Error("Release did not drop the node tree")
	}

	// Now the document's own turn to be tested.
	doc.Parse()

	if doc.BasePath != basePath {
		t.Errorf("BasePath changed on re-parse: %q != %q", doc.BasePath, basePath)
	}
	if len(doc.NodesOfInterest) != nodes {
		t.Errorf("NodesOfInterest changed on re-parse: %d != %d",
			len(doc.NodesOfInterest), nodes)
	}
	if len(doc.hashMap) != hashes {
		t.Errorf("hashMap changed on re-parse: %d != %d", len(doc.hashMap), hashes)
	}
	if doc.htmlNode == nil {
		t.Error("re-parse did not reload the node tree")
	}
}
