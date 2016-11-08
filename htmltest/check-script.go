package htmltest

import (
	// "fmt"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
)

func CheckScript(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"src", "data-proofer-ignore"})

	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, "data-proofer-ignore") {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])
	_ = ref

	// Check src problems
	if attrPresent(node.Attr, "src") && len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute present but empty",
			Reference: ref,
		})
		return
	}

	// Check invalid content
	if !attrPresent(node.Attr, "src") && node.FirstChild == nil {
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "script content missing / no src attribute",
			Reference: ref,
		})
		return
	}

	// Route reference check
	switch ref.Scheme() {
	case "http":
		if Opts.EnforceHTTPS {
			issues.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "is not an HTTPS target",
				Reference: ref,
			})
		}
		CheckExternal(ref)
	case "https":
		CheckExternal(ref)
	case "file":
		CheckInternal(ref)
	}
}
