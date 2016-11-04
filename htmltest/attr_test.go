package htmltest

import (
	// "fmt"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestExtractAttrs(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild
	attrs := extractAttrs(nodeImg.Attr, []string{"src", "alt"})

	t_assertEqual(t, attrs["src"], "x")
	t_assertEqual(t, attrs["alt"], "y")
	t_assertNotEqual(t, attrs["foo"], "bar")
}

func TestAttrPresent(t *testing.T) {
	snip := "<img src=\"x\" alt=\"y\" />"
	nodeDoc, _ := html.Parse(strings.NewReader(snip))
	nodeImg := nodeDoc.FirstChild.FirstChild.NextSibling.FirstChild

	t_assertEqual(t, attrPresent(nodeImg.Attr, "src"), true)
	t_assertEqual(t, attrPresent(nodeImg.Attr, "alt"), true)
	t_assertNotEqual(t, attrPresent(nodeImg.Attr, "foo"), true)
}
