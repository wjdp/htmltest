package issues

import (
	"github.com/daviddengcn/go-assert"
	"github.com/wjdp/htmltest/htmldoc"
	"testing"
)

func TestIssuePrimary(t *testing.T) {
	issue0 := Issue{}
	assert.Equals(t, "issue0 primary", issue0.primary(), TEXT_NIL)

	doc := htmldoc.Document{
		SitePath: "dir/doc.html",
	}
	issue1 := Issue{
		Document: &doc,
	}
	assert.Equals(t, "issue1 primary", issue1.primary(), "dir/doc.html")

	ref := htmldoc.Reference{
		Document: &doc,
	}
	issue2 := Issue{
		Reference: &ref,
	}
	assert.Equals(t, "issue2 primary", issue2.primary(), "dir/doc.html")
}

func TestIssueSecondary(t *testing.T) {
	issue0 := Issue{}
	assert.Equals(t, "issue0 secondary", issue0.secondary(), TEXT_NIL)

	ref := htmldoc.Reference{
		Path: "http://example.com",
	}
	issue1 := Issue{
		Reference: &ref,
	}
	assert.Equals(t, "issue1 secondary", issue1.secondary(), "http://example.com")
}

func ExampleIssuePrint() {
	doc := htmldoc.Document{
		SitePath: "dir/doc.html",
	}
	ref := htmldoc.Reference{
		Document: &doc,
		Path:     "http://example.com",
	}

	issueStore := IssueStore{
		LogLevel: WARNING,
	}

	issue1 := Issue{
		Level:    ERROR,
		Document: &doc,
		store:    &issueStore,
		Message:  "test1",
	}
	issue1.print(false)

	issue2 := Issue{
		Level:     WARNING,
		Reference: &ref,
		store:     &issueStore,
		Message:   "test2",
	}
	issue2.print(false)

	issue3 := Issue{
		Level:    INFO,
		Document: &doc,
		store:    &issueStore,
		Message:  "test3",
	}
	issue3.print(false)

	// Output:
	// test1 --- dir/doc.html --> <nil>
	// test2 --- dir/doc.html --> http://example.com

}
