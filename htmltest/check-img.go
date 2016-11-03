package htmltest

import (
	// "log"
	"github.com/wjdp/htmltest/htmldoc"
	"golang.org/x/net/html"
)

func CheckImg(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"href", "alt", "data-proofer-ignore"})

	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, "data-proofer-ignore") {
		return
	}
	_ = attrs
}
