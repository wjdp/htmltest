package htmltest

import (
	"fmt"
	"github.com/theunrepentantgeek/htmltest/htmldoc"
	"github.com/theunrepentantgeek/htmltest/issues"
	"golang.org/x/net/html"
	"regexp"
)

func (hT *HTMLTest) checkMeta(document *htmldoc.Document, node *html.Node) {
	if hT.opts.CheckMetaRefresh {
		hT.checkMetaRefresh(document, node)
	}
}

func (hT *HTMLTest) checkMetaRefresh(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"http-equiv", "content"})

	// Checks for meta refresh redirect tag
	if attrs["http-equiv"] == "refresh" {
		// Extract the timing and path from the content attr

		// Build regex to match ;url= and alike, the split the content attribute
		re, _ := regexp.Compile(";[ ]{0,1}[Uu][Rr][Ll]=")
		contentSplit := re.Split(attrs["content"], 2)

		// Define ref from this
		var ref *htmldoc.Reference
		var err error
		if len(contentSplit) == 2 {
			if contentSplit[1][0] == 34 || contentSplit[1][0] == 39 {
				hT.issueStore.AddIssue(issues.Issue{
					Level:    issues.LevelError,
					Message:  "url in meta refresh must not start with single or double quote",
					Document: document,
				})
				return
			}
			ref, err = htmldoc.NewReference(document, node, contentSplit[1])
		} else {
			ref, err = htmldoc.NewReference(document, node, "")
		}
		if err != nil {
			hT.issueStore.AddIssue(issues.Issue{
				Level:    issues.LevelError,
				Document: document,
				Message:  fmt.Sprintf("bad reference: %q", err),
			})
			return
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
