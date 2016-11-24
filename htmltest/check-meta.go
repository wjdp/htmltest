package htmltest

import (
	// "fmt"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func (hT *HTMLTest) checkMeta(document *htmldoc.Document, node *html.Node) {
	hT.checkMetaRefresh(document, node)
}

func (hT *HTMLTest) checkMetaRefresh(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"http-equiv", "content", hT.opts.IgnoreTagAttribute})

	// Checks for meta refresh redirect tag
	if attrs["http-equiv"] == "refresh" {
		// Extract the timing and path from the content attr
		contentSplit := strings.Split(attrs["content"], ";url=")
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
