package htmltest

import (
	"fmt"

	"github.com/theunrepentantgeek/htmltest/htmldoc"
	"github.com/theunrepentantgeek/htmltest/issues"
	"golang.org/x/net/html"
)

// Checks the reference in the provided node and attribute key
func (hT *HTMLTest) checkGeneric(document *htmldoc.Document, node *html.Node, key string) {
	// Fail silently if attribute isn't present
	if !htmldoc.AttrPresent(node.Attr, key) {
		return
	}

	urlStr := htmldoc.GetAttr(node.Attr, key)
	ref, err := htmldoc.NewReference(document, node, urlStr)
	if err != nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:    issues.LevelError,
			Document: document,
			Message:  fmt.Sprintf("bad reference: %q", err),
		})
		return
	}

	// Check attr isn't blank
	if urlStr == "" {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   fmt.Sprintf(node.Data, key, "is blank"),
			Reference: ref,
		})
	}

	// Check the reference
	hT.checkGenericRef(ref)
}

func (hT *HTMLTest) checkGenericRef(ref *htmldoc.Reference) {
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

func (hT *HTMLTest) enforceHTTPS(ref *htmldoc.Reference) {
	urlStr := ref.URLString()

	// Does this url match an url ignore rule?
	if hT.opts.isURLIgnored(urlStr) || hT.opts.isInsecureURLIgnored(urlStr) {
		return
	}
	issueLevel := issues.LevelError
	if hT.opts.IgnoreExternalBrokenLinks {
		issueLevel = issues.LevelWarning
	}

	if hT.opts.EnforceHTTPS {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issueLevel,
			Message:   "is not an HTTPS target",
			Reference: ref,
		})
	}
}
