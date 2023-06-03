package issues

import (
	"fmt"
	"github.com/daviddengcn/go-assert"
	"github.com/wjdp/htmltest/htmldoc"
	"io/ioutil"
	"os"
	"reflect"
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

func TestGetIssueStats_None(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	stats := iS.GetIssueStats()
	assert.IsTrue(t, "TotalByLevel", reflect.DeepEqual(stats.TotalByLevel, map[int]int{}))
	assert.IsTrue(t, "ErrorsByMessage", reflect.DeepEqual(stats.ErrorsByMessage, map[string]int{}))
	assert.IsTrue(t, "WarningsByMessage", reflect.DeepEqual(stats.WarningsByMessage, map[string]int{}))
}

func addOneError(iS *IssueStore) {
	iS.AddIssue(Issue{
		Level:   LevelError,
		Message: "test",
	})
}

func addMultipleIssues(iS *IssueStore) {
	iS.AddIssue(Issue{
		Level:   LevelError,
		Message: "test1",
	})
	iS.AddIssue(Issue{
		Level:   LevelWarning,
		Message: "test1",
	})
	iS.AddIssue(Issue{
		Level:   LevelInfo,
		Message: "test1",
	})
	iS.AddIssue(Issue{
		Level:   LevelDebug,
		Message: "test1",
	})
	iS.AddIssue(Issue{
		Level:   LevelError,
		Message: "test2",
	})
	iS.AddIssue(Issue{
		Level:   LevelError,
		Message: "test2",
	})
}

func TestGetIssueStats_OneError(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	addOneError(&iS)
	stats := iS.GetIssueStats()
	assert.Equals(
		t, "TotalByLevel",
		fmt.Sprint(stats.TotalByLevel),
		fmt.Sprint(map[int]int{LevelError: 1}),
	)
	assert.Equals(
		t, "ErrorsByMessage",
		fmt.Sprint(stats.ErrorsByMessage),
		fmt.Sprint(map[string]int{"test": 1}),
	)
}

func TestGetIssueStats_MultipleIssues(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	addMultipleIssues(&iS)
	stats := iS.GetIssueStats()
	assert.Equals(t,
		"TotalByLevel",
		fmt.Sprint(stats.TotalByLevel),
		fmt.Sprint(map[int]int{LevelError: 3, LevelWarning: 1, LevelInfo: 1, LevelDebug: 1}))
	assert.Equals(
		t, "ErrorsByMessage",
		fmt.Sprint(stats.ErrorsByMessage),
		fmt.Sprint(map[string]int{"test1": 1, "test2": 2}),
	)
	assert.Equals(
		t, "WarningsByMessage",
		fmt.Sprint(stats.WarningsByMessage),
		fmt.Sprint(map[string]int{"test1": 1}),
	)
}

func TestFormatIssueStats_None(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	assert.Equals(t, "FormatIssueStats", iS.FormatIssueStats(), "  Errors:   0\n")
}

func TestFormatIssueStats_Multiple_AtLogLevelError(t *testing.T) {
	iS := NewIssueStore(LevelError, false)
	addMultipleIssues(&iS)
	const expected = `  Errors:   3
  Errors by message:
    1 "test1"
    2 "test2"
`
	assert.Equals(t, "FormatIssueStats", iS.FormatIssueStats(), expected)
}

func TestFormatIssueStats_Multiple_AtLogLevelWarning(t *testing.T) {
	iS := NewIssueStore(LevelWarning, false)
	addMultipleIssues(&iS)
	const expected = `  Errors:   3
  Warnings: 1
  Errors by message:
    1 "test1"
    2 "test2"
  Warnings by message:
    1 "test1"
`
	assert.Equals(t, "FormatIssueStats", iS.FormatIssueStats(), expected)
}
