package htmltest

import (
	"golang.org/x/net/html"
	"regexp"
	"wjdp.uk/htmltest/htmldoc"
	"wjdp.uk/htmltest/issues"
)

func (hT *HTMLTest) checkMeta(document *htmldoc.Document, node *html.Node) {
	if hT.opts.CheckMetaRefresh {
		hT.checkMetaRefresh(document, node)
	}
}

func (hT *HTMLTest) checkMetaRefresh(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"http-equiv", "content", hT.opts.IgnoreTagAttribute})

	// Checks for meta refresh redirect tag
	if attrs["http-equiv"] == "refresh" {
		// Extract the timing and path from the content attr

		// Build regex to match ;url= and alike, the split the content attribute
		re, _ := regexp.Compile(";[ ]{0,1}[Uu][Rr][Ll]=")
		contentSplit := re.Split(attrs["content"], 2)

		// Define ref from this
		var ref *htmldoc.Reference
		if len(contentSplit) == 2 {
			ref = htmldoc.NewReference(document, node, contentSplit[1])
		} else {
			ref = htmldoc.NewReference(document, node, "")
		}

		// If refresh the content attribute must be set
		if htmldoc.AttrPresent(node.Attr, "content") {

			// Check content isn't blank
			if len(attrs["content"]) == 0 {
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.LevelError,
					Message:   "blank content attribute in meta refresh",
					Reference: ref,
				})
				return // stop
			}

			// Check the time is a positive integer, if the user has buggered up
			// with ;url=... this will also display.
			if ok, _ := regexp.MatchString("^\\d+$", contentSplit[0]); !ok {
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.LevelError,
					Message:   "invalid content attribute in meta refresh",
					Reference: ref,
				})
			}

			// Check the reference
			hT.checkGenericRef(ref)

		} else {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "missing content attribute in meta refresh",
				Reference: ref,
			})
		}

	}
}
