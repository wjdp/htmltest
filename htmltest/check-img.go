package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"regexp"
)

func (hT *HTMLTest) checkImg(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"src", "alt", "usemap", hT.opts.IgnoreTagAttribute})

	// Ignore if data-proofer-ignore set
	if htmldoc.AttrPresent(node.Attr, hT.opts.IgnoreTagAttribute) {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])

	// Check alt present, fail if absent unless asked to ignore
	if !htmldoc.AttrPresent(node.Attr, "alt") && !hT.opts.IgnoreAltMissing {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "alt attribute missing",
			Reference: ref,
		})
	} else if htmldoc.AttrPresent(node.Attr, "alt") {
		// Following checks require alt attr is present
		if len(attrs["alt"]) == 0 {
			// Check alt has length, fail if empty
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "alt text empty",
				Reference: ref,
			})
		}
		if b, _ := regexp.MatchString("^\\s+$", attrs["alt"]); b {
			// Check alt is not just whitespace
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "alt text contains only whitespace",
				Reference: ref,
			})
		}
	}

	// Check src present, fail if absent
	if !htmldoc.AttrPresent(node.Attr, "src") {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "src attribute missing",
			Reference: ref,
		})
		return
	} else if len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "src attribute empty",
			Reference: ref,
		})
		return
	}

	// Check usemap
	if htmldoc.AttrPresent(node.Attr, "usemap") {
		usemapRef := htmldoc.NewReference(document, node, attrs["usemap"])

		if len(usemapRef.URL.Path) > 0 {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "only fragment starting with # allowed in usemap attribute",
				Reference: ref,
			})
		} else if len(usemapRef.URL.Fragment) == 0 {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "usemap empty",
				Reference: ref,
			})
		} else {
			hT.checkInternalHash(usemapRef)
		}

		parent := node.Parent
		if parent.Data == "a" {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "<img> with usemap attribute not allowed as descendant of an <a> element",
				Reference: ref,
			})
		} else if parent.Data == "button" {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "<img> with usemap attribute not allowed as descendant of a <button>",
				Reference: ref,
			})
		}
	}

	// Route reference check
	switch ref.Scheme() {
	case "http":
		hT.enforceHTTPS(ref)
		hT.checkExternal(ref)
	case "https":
		hT.checkExternal(ref)
	case "file":
		hT.checkInternal(ref)
	}
}
