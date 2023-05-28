package issues

import (
	"github.com/daviddengcn/go-assert"
	"github.com/theunrepentantgeek/htmltest/htmldoc"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestIssueStoreNew(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	assert.Equals(t, "IssueStore LogLevel", iS.logLevel, LevelError)
}

func TestIssueStoreAdd(t *testing.T) {
	iS := NewIssueStore(LevelNone, false)
	issue := Issue{Level: LevelError, Message: "test"}
	iS.AddIssue(issue)
	assert.Equals(t, "issue count", iS.Count(LevelError), 1)
}

func TestIssueStoreCount(t *testing.T) {
	iS := NewIssueStore(LevelNone, false)
	iS.AddIssue(Issue{Level: LevelError, Message: "error"})
	iS.AddIssue(Issue{Level: LevelWarning, Message: "warn"})
	iS.AddIssue(Issue{Level: LevelInfo, Message: "notice"})
	assert.Equals(t, "issue count", iS.Count(LevelError), 1)
	assert.Equals(t, "issue count", iS.Count(LevelWarning), 2)
	assert.Equals(t, "issue count", iS.Count(LevelInfo), 3)
}

func TestIssueStoreMessageMatchCount(t *testing.T) {
	iS := NewIssueStore(LevelNone, false)
	iS.AddIssue(Issue{Level: LevelError, Message: "error one"})
	iS.AddIssue(Issue{Level: LevelWarning, Message: "error two"})
	iS.AddIssue(Issue{Level: LevelInfo, Message: "notice"})
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
	iS := NewIssueStore(LevelError, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue1 := Issue{
		Level:    LevelError,
		Message:  "test1",
		Document: &doc,
	}
	iS.AddIssue(issue1)
	issue2 := Issue{
		Level:    LevelWarning,
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
	iS := NewIssueStore(LevelNone, true)
	issue1 := Issue{
		Level:   LevelError,
		Message: "test1",
	}
	iS.AddIssue(issue1)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue2 := Issue{
		Level:    LevelError,
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
	iS := NewIssueStore(LevelError, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue := Issue{
		Level:    LevelError,
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
	iS := NewIssueStore(LevelError, false)
	doc := htmldoc.Document{
		SitePath: "dir/page.html",
	}
	issue := Issue{
		Level:    LevelInfo,
		Message:  "test1",
		Document: &doc,
	}
	iS.AddIssue(issue)

	iS.PrintDocumentIssues(&doc)
	// Output:
}
