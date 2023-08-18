package htmltest

import (
	"fmt"

	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
)

func (hT *HTMLTest) checkDoctype(document *htmldoc.Document) {
	// Error if no doctype
	// The doctype *must* be the first element in the document
	// If it's not golang.org/x/net/html simply ignores it.
	if document.DoctypeNode == nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:    issues.LevelError,
			Message:  "missing doctype",
			Document: document,
		})
		return
	}

	// Dump the doctype data and attrs to debug
	hT.issueStore.AddIssue(issues.Issue{
		Level: issues.LevelDebug,
		Message: fmt.Sprintf("DOCTYPE %+v %+v\n",
			document.DoctypeNode.Data, document.DoctypeNode.Attr),
		Document: document,
	})

	isHTML5 := (document.DoctypeNode.Data == "html" &&
		len(document.DoctypeNode.Attr) == 0)

	if hT.opts.EnforceHTML5 && !isHTML5 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:    issues.LevelError,
			Message:  "doctype isn't html5",
			Document: document,
		})
	}

}
