package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"regexp"
)

func (hT *HtmlTest) checkImg(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr,
		[]string{"src", "alt", hT.opts.IgnoreTagAttribute})

	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, hT.opts.IgnoreTagAttribute) {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])

	// Check alt present, fail if absent unless asked to ignore
	if !attrPresent(node.Attr, "alt") && !hT.opts.IgnoreAltMissing {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "alt attribute missing",
			Reference: ref,
		})
	} else if attrPresent(node.Attr, "alt") {
		// Following checks require alt attr is present
		if len(attrs["alt"]) == 0 {
			// Check alt has length, fail if empty
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "alt text empty",
				Reference: ref,
			})
		}
		if b, _ := regexp.MatchString("^\\s+$", attrs["alt"]); b {
			// Check alt is not just whitespace
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "alt text contains only whitespace",
				Reference: ref,
			})
		}
	}

	// Check src present, fail if absent
	if !attrPresent(node.Attr, "src") {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute missing",
			Reference: ref,
		})
		return
	} else if len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute empty",
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
