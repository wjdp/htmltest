package htmltest

import (
	"golang.org/x/net/html"
	"wjdp.uk/htmltest/htmldoc"
	"wjdp.uk/htmltest/issues"
)

func (hT *HTMLTest) checkScript(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"src", hT.opts.IgnoreTagAttribute})

	// Ignore if data-proofer-ignore set
	if htmldoc.AttrPresent(node.Attr, hT.opts.IgnoreTagAttribute) {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["src"])

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
