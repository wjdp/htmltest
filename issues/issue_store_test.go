package issues

import (
	"github.com/daviddengcn/go-assert"
	"github.com/wjdp/htmltest/htmldoc"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestIssueStoreNew(t *testing.T) {
	iS := NewIssueStore(ERROR, false)
	assert.Equals(t, "IssueStore LogLevel", iS.logLevel, ERROR)
}

func TestIssueStoreAdd(t *testing.T) {
	iS := NewIssueStore(NONE, false)
	issue := Issue{Level: ERROR, Message: "test"}
	iS.AddIssue(issue)
	assert.Equals(t, "issue count", iS.Count(ERROR), 1)
}

func TestIssueStoreCount(t *testing.T) {
	iS := NewIssueStore(NONE, false)
	iS.AddIssue(Issue{Level: ERROR, Message: "error"})
	iS.AddIssue(Issue{Level: WARNING, Message: "warn"})
	iS.AddIssue(Issue{Level: INFO, Message: "notice"})
	assert.Equals(t, "issue count", iS.Count(ERROR), 1)
	assert.Equals(t, "issue count", iS.Count(WARNING), 2)
	assert.Equals(t, "issue count", iS.Count(INFO), 3)
}

func TestIssueStoreMessageMatchCount(t *testing.T) {
	iS := NewIssueStore(NONE, false)
	iS.AddIssue(Issue{Level: ERROR, Message: "error one"})
	iS.AddIssue(Issue{Level: WARNING, Message: "error two"})
	iS.AddIssue(Issue{Level: INFO, Message: "notice"})
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("carrots"), 0)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("error"), 2)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("two"), 1)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("notice"), 1)
}

func TestIssueStoreWriteLog(t *testing.T) {
	// passes for log written using LogLevel
	LOGFILE := "issue-store-test.log"
	iS := NewIssueStore(ERROR, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue1 := Issue{
		Level:    ERROR,
		Message:  "test1",
		Document: &doc,
	}
	iS.AddIssue(issue1)
	issue2 := Issue{
		Level:    WARNING,
		Message:  "test2",
		Document: &doc,
	}
	iS.AddIssue(issue2)

	iS.WriteLog(LOGFILE)
	logBytes, err := ioutil.ReadFile(LOGFILE)
	assert.Equals(t, "file error", err, nil)
	logString := string(logBytes)

	assert.IsTrue(t, "log contents", strings.Contains(
		logString, "test1 --- dir/page.html --> <nil>"))
	assert.IsFalse(t, "log contents", strings.Contains(
		logString, "test2 --- dir/page.html --> <nil>"))

	removeErr := os.Remove(LOGFILE)
	assert.Equals(t, "file error", removeErr, nil)

}

func ExampleIssueStoreDumpIssues() {
	// Passes for dumping all issues, ignoring LogLevel
	iS := NewIssueStore(NONE, true)
	issue1 := Issue{
		Level:   ERROR,
		Message: "test1",
	}
	iS.AddIssue(issue1)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue2 := Issue{
		Level:    ERROR,
		Message:  "test2",
		Document: &doc,
	}
	iS.AddIssue(issue2)

	iS.DumpIssues(true)
	// Output:
	// <<<<<<<<<<<<<<<<<<<<<<<<
	// test1
	// test2 --- dir/page.html --> <nil>
	// >>>>>>>>>>>>>>>>>>>>>>>>
}

func ExampleIssueStorePrintDocumentIssues() {
	iS := NewIssueStore(ERROR, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue := Issue{
		Level:    ERROR,
		Message:  "test1",
		Document: &doc,
	}
	iS.AddIssue(issue)

	iS.PrintDocumentIssues(&doc)
	// Output:
	// dir/page.html
	//   test1 --- dir/page.html --> <nil>
}

func ExampleIssueStorePrintDocumentIssuesEmpty() {
	iS := NewIssueStore(ERROR, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue := Issue{
		Level:    INFO,
		Message:  "test1",
		Document: &doc,
	}
	iS.AddIssue(issue)

	iS.PrintDocumentIssues(&doc)
	// Output:
}
