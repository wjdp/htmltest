package htmltest

import (
	// "fmt"
	"github.com/wjdp/htmltest/htmldoc"
	// "github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
)

func CheckScript(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"src", "alt", "data-proofer-ignore"})
	ref := htmldoc.NewReference(document, node, attrs["src"])
	_ = ref
}
