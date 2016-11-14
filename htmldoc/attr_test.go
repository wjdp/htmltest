package htmldoc

import (
	"github.com/daviddengcn/go-assert"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestGetAttr(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	assert.Equals(t, "src", GetAttr(nodeImg.Attr, "src"), "x")
	assert.Equals(t, "alt", GetAttr(nodeImg.Attr, "alt"), "y")
}

func TestExtractAttrs(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild
	attrs := ExtractAttrs(nodeImg.Attr, []string{"src", "alt"})

	assert.Equals(t, "src", attrs["src"], "x")
	assert.Equals(t, "alt", attrs["alt"], "y")
	assert.NotEquals(t, "foo", attrs["foo"], "bar")
}

func TestAttrPresent(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	assert.Equals(t, "src in attr", AttrPresent(nodeImg.Attr, "src"), true)
	assert.Equals(t, "alt in attr", AttrPresent(nodeImg.Attr, "src"), true)
	assert.NotEquals(t, "foo in attr", AttrPresent(nodeImg.Attr, "src"), false)
}

func TestAttrValIdId(t *testing.T) {
	snip := "<h1 id=\"x\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeH1 := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	assert.Equals(t, "h1 id", GetId(nodeH1.Attr), "x")
}

func TestAttrValIdName(t *testing.T) {
	snip := "<h1 name=\"x\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeH1 := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	assert.Equals(t, "h1 name", GetId(nodeH1.Attr), "x")
}
