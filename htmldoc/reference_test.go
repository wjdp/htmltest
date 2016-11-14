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
		SitePath: "doc.html",
		htmlNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "http://test.com")
	assert.Equals(t, "http reference", ref.Scheme(), "http")
	ref = NewReference(&doc, nodeElem, "https://test.com")
	assert.Equals(t, "https reference", ref.Scheme(), "https")
	ref = NewReference(&doc, nodeElem, "//test.com")
	assert.Equals(t, "https reference", ref.Scheme(), "https")
	ref = NewReference(&doc, nodeElem,
		"https://photos.smugmug.com/photos/i-CNHsHLM/0/440x622/i-CNHsHLM-440x622.jpg")
	assert.Equals(t, "http reference", ref.Scheme(), "https")
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

	// Grubby url
	ref = NewReference(&doc, nodeElem, "\n http://foo")
	assert.Equals(t, "unknown reference", ref.Scheme(), "http")
}

func TestReferenceURLString(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		SitePath: "doc.html",
		htmlNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "http://example.com")
	assert.Equals(t, "URLString", ref.URLString(), "http://example.com")
	ref = NewReference(&doc, nodeElem, "http://example.com/")
	assert.Equals(t, "URLString", ref.URLString(), "http://example.com/")
	ref = NewReference(&doc, nodeElem, "https://example.com")
	assert.Equals(t, "URLString", ref.URLString(), "https://example.com")
	ref = NewReference(&doc, nodeElem, "//example.com")
	assert.Equals(t, "URLString", ref.URLString(), "https://example.com")

}

func TestReferenceIsInternalAbsolute(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, nodeElem := nodeGen(snip)

	doc := Document{
		SitePath: "doc.html",
		htmlNode: nodeDoc,
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
		SitePath: "doc.html",
		BasePath: "directory/subdir",
		htmlNode: nodeDoc,
	}

	var ref *Reference

	ref = NewReference(&doc, nodeElem, "/abc/page.html")
	assert.Equals(t, "internal absolute reference", ref.RefSitePath(), "/abc/page.html")
	ref = NewReference(&doc, nodeElem, "/yyz")
	assert.Equals(t, "internal absolute reference", ref.RefSitePath(), "/yyz")
	ref = NewReference(&doc, nodeElem, "zzy")
	assert.Equals(t, "internal relative reference", ref.RefSitePath(), "directory/subdir/zzy")
	ref = NewReference(&doc, nodeElem, "zzy/uup.jjr")
	assert.Equals(t, "internal relative reference", ref.RefSitePath(), "directory/subdir/zzy/uup.jjr")
	ref = NewReference(&doc, nodeElem, "./zzy/uup.jjr")
	assert.Equals(t, "internal relative reference", ref.RefSitePath(), "directory/subdir/zzy/uup.jjr")
}

func TestURLStripQueryString(t *testing.T) {
	original := "https://github.com/wjdp/gotdict/issues/new?title=Harwood Fell&body=[_definitions/harwood-fell.mdd](https://github.com/wjdp/gotdict/blob/master/_definitions/harwood-fell.mdd)"
	actual := URLStripQueryString(original)
	expected := "https://github.com/wjdp/gotdict/issues/new"

	assert.Equals(t, "stripped url", actual, expected)
}
