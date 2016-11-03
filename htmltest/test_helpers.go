package htmltest

import (
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

const t_LogLevel int = issues.WARNING

func t_expectIssue(t *testing.T, message string, expected int) {
	c := issues.IssueMatchCount(message)
	if c != expected {
		t.Error("expected issue", message, "count", expected, "!=", c)
		issues.OutputIssues()
	}
}

func t_expectIssueCount(t *testing.T, expected int) {
	c := issues.IssueCount(issues.WARNING)
	if c != expected {
		t.Error("expected", expected, "issues,", c, "found")
		issues.OutputIssues()
	}
}

func t_testFile(filename string) {
	opts := map[string]interface{}{
		"DirectoryPath": path.Dir(filename),
		"FilePath":      path.Base(filename),
		"LogLevel":      t_LogLevel,
	}
	Test(opts)
}

func t_testDirectory(filename string) {
	opts := map[string]interface{}{
		"DirectoryPath": path.Dir(filename),
		"LogLevel":      t_LogLevel,
	}
	Test(opts)
}
