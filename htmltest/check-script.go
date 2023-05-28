package htmltest

import (
	"fmt"
	"github.com/theunrepentantgeek/htmltest/htmldoc"
	"github.com/theunrepentantgeek/htmltest/issues"
	"golang.org/x/net/html"
)

func (hT *HTMLTest) checkScript(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"src"})

	// Create reference
	ref, err := htmldoc.NewReference(document, node, attrs["src"])
	if err != nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:    issues.LevelError,
			Document: document,
			Message:  fmt.Sprintf("bad reference: %q", err),
		})
		return
	}

	// Check src problems
	if htmldoc.AttrPresent(node.Attr, "src") && len(attrs["src"]) == 0 {
		// Check src has length, fail if empty
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "src attribute present but empty",
			Reference: ref,
		})
		return
	}

	// Check invalid content
	if !htmldoc.AttrPresent(node.Attr, "src") && node.FirstChild == nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "script content missing / no src attribute",
			Reference: ref,
		})
		return
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
