package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"regexp"
)

func CheckImg(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"src", "alt", "data-proofer-ignore"})

	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, "data-proofer-ignore") {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])

	// Check src present, fail if absent
	if !attrPresent(node.Attr, "src") {
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute missing",
			Reference: ref,
		})
		return
	} else if len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "src attribute empty",
			Reference: ref,
		})
		return
	}

	// Check alt present, fail if absent unless asked to ignore
	if !attrPresent(node.Attr, "alt") && !Opts.IgnoreAlt {
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "alt attribute missing",
			Reference: ref,
		})
		return
	} else if len(attrs["alt"]) == 0 && !Opts.IgnoreAlt {
		// Check alt has length, fail if empty unless asked to ignore
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "alt text empty",
			Reference: ref,
		})
		return
	} else if b, _ := regexp.MatchString("^\\s+$", attrs["alt"]); b {
		// Check alt is not just whitespace
		issues.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "alt text contains only whitespace",
			Reference: ref,
		})
	}

	// Route reference check
	switch ref.Scheme {
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
