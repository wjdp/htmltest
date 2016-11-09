package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
)

func (hT *HtmlTest) checkScript(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"src", "data-proofer-ignore"})

	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, "data-proofer-ignore") {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])

	// Check src problems
	if attrPresent(node.Attr, "src") && len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute present but empty",
			Reference: ref,
		})
		return
	}

	// Check invalid content
	if !attrPresent(node.Attr, "src") && node.FirstChild == nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "script content missing / no src attribute",
			Reference: ref,
		})
		return
	}

	// Route reference check
	switch ref.Scheme() {
	case "http":
		if hT.opts.EnforceHTTPS {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "is not an HTTPS target",
				Reference: ref,
			})
		}
		hT.checkExternal(ref)
	case "https":
		hT.checkExternal(ref)
	case "file":
		hT.checkInternal(ref)
	}
}
