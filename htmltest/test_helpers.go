package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

func t_assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Error(a, "!=", b)
	}
}

func t_assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Error(a, "==", b)
	}
}

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

func t_testFileOpts(filename string, t_opts map[string]interface{}) {
	opts := map[string]interface{}{
		"DirectoryPath": path.Dir(filename),
		"FilePath":      path.Base(filename),
		"LogLevel":      t_LogLevel,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	Test(opts)
}

func t_testDirectory(filename string) {
	opts := map[string]interface{}{
		"DirectoryPath": filename,
		"LogLevel":      t_LogLevel,
	}
	Test(opts)
}

func t_testDirectoryOpts(filename string, t_opts map[string]interface{}) {
	opts := map[string]interface{}{
		"DirectoryPath": filename,
		"LogLevel":      t_LogLevel,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	Test(opts)
}
