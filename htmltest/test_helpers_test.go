package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

const t_LogLevel int = issues.NONE

func t_expectIssue(t *testing.T, hT *HtmlTest, message string, expected int) {
	c := hT.issueStore.MessageMatchCount(message)
	if c != expected {
		t.Error("expected issue", message, "count", expected, "!=", c)
	}
}

func t_expectIssueCount(t *testing.T, hT *HtmlTest, expected int) {
	c := hT.issueStore.Count(issues.WARNING)
	if c != expected {
		t.Error("expected", expected, "issues,", c, "found")
	}
}

func t_testFile(filename string) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath": path.Dir(filename),
		"FilePath":      path.Base(filename),
		"LogLevel":      t_LogLevel,
	}
	return Test(opts)
}

func t_testFileOpts(filename string, t_opts map[string]interface{}) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath": path.Dir(filename),
		"FilePath":      path.Base(filename),
		"LogLevel":      t_LogLevel,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	return Test(opts)
}

func t_testDirectory(filename string) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath": filename,
		"LogLevel":      t_LogLevel,
	}
	return Test(opts)
}

func t_testDirectoryOpts(filename string, t_opts map[string]interface{}) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath": filename,
		"LogLevel":      t_LogLevel,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	return Test(opts)
}
