package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func nodeGen(snip string) (*html.Node, *html.Node) {
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeElem := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild
	return nodeDoc, nodeElem
}

func TestReferenceScheme(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		Path:     "doc.html",
		HTMLNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "http://test.com")
	assert.Equals(t, "http reference", ref.Scheme(), "http")
	ref = NewReference(&doc, nodeElem, "https://test.com")
	assert.Equals(t, "https reference", ref.Scheme(), "https")

	ref = NewReference(&doc, nodeElem, "x?a=1#3")
	assert.Equals(t, "file reference", ref.Scheme(), "file")
	ref = NewReference(&doc, nodeElem, "#123")
	assert.Equals(t, "self reference", ref.Scheme(), "self")
	ref = NewReference(&doc, nodeElem, "mailto:x@y.com")
	assert.Equals(t, "mailto reference", ref.Scheme(), "mailto")
	ref = NewReference(&doc, nodeElem, "tel:123")
	assert.Equals(t, "tel reference", ref.Scheme(), "tel")
	ref = NewReference(&doc, nodeElem, "abc:123")
	assert.Equals(t, "unknown reference", ref.Scheme(), "")
}

func TestReferenceURLString(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		Path:     "doc.html",
		HTMLNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "google.com")
	assert.Equals(t, "URLString", ref.URLString(), "google.com")

}

func TestReferenceIsInternalAbsolute(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		Path:     "doc.html",
		HTMLNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "/abc/page.html")
	assert.IsTrue(t, "internal absolute reference", ref.IsInternalAbsolute())
	ref = NewReference(&doc, nodeElem, "/yyz")
	assert.IsTrue(t, "internal absolute reference", ref.IsInternalAbsolute())
	ref = NewReference(&doc, nodeElem, "zzy")
	assert.IsFalse(t, "internal relative reference", ref.IsInternalAbsolute())
	ref = NewReference(&doc, nodeElem, "zzy/uup.jjr")
	assert.IsFalse(t, "internal relative reference", ref.IsInternalAbsolute())
	ref = NewReference(&doc, nodeElem, "./zzy/uup.jjr")
	assert.IsFalse(t, "internal relative reference", ref.IsInternalAbsolute())
}

func TestReferenceAbsolutePath(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		Path:      "doc.html",
		Directory: "directory/subdir",
		HTMLNode:  nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "/abc/page.html")
	assert.Equals(t, "internal absolute reference", ref.AbsolutePath(), "/abc/page.html")
	ref = NewReference(&doc, nodeElem, "/yyz")
	assert.Equals(t, "internal absolute reference", ref.AbsolutePath(), "/yyz")
	ref = NewReference(&doc, nodeElem, "zzy")
	assert.Equals(t, "internal relative reference", ref.AbsolutePath(), "directory/subdir/zzy")
	ref = NewReference(&doc, nodeElem, "zzy/uup.jjr")
	assert.Equals(t, "internal relative reference", ref.AbsolutePath(), "directory/subdir/zzy/uup.jjr")
	ref = NewReference(&doc, nodeElem, "./zzy/uup.jjr")
	assert.Equals(t, "internal relative reference", ref.AbsolutePath(), "directory/subdir/zzy/uup.jjr")
}
