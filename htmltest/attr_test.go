package htmltest

import (
	"github.com/daviddengcn/go-assert"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestExtractAttrs(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild
	attrs := extractAttrs(nodeImg.Attr, []string{"src", "alt"})

	assert.Equals(t, "src", attrs["src"], "x")
	assert.Equals(t, "alt", attrs["alt"], "y")
	assert.NotEquals(t, "foo", attrs["foo"], "bar")
}

func TestAttrPresent(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	assert.Equals(t, "src in attr", attrPresent(nodeImg.Attr, "src"), true)
	assert.Equals(t, "alt in attr", attrPresent(nodeImg.Attr, "src"), true)
	assert.NotEquals(t, "foo in attr", attrPresent(nodeImg.Attr, "src"), false)
}
