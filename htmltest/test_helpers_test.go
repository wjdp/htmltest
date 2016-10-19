package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

// Kepe it quiet
const t_LogLevel int = issues.NONE

// We're running non-concurrently, speed up the tests by turning down the
// timeout. Assumes we're on a good connection.
const t_ExternalTimeout int = 3

func t_expectIssue(t *testing.T, hT *HtmlTest, message string, expected int) {
	c := hT.issueStore.MessageMatchCount(message)
	if c != expected {
		t.Error("expected issue", message, "count", expected, "!=", c)
	}
}

func t_expectIssueCount(t *testing.T, hT *HtmlTest, expected int) {
	c := hT.issueStore.Count(issues.ERROR)
	if c != expected {
		t.Error("expected", expected, "issues,", c, "found")
	}
}

func t_testFile(filename string) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath":   path.Dir(filename),
		"FilePath":        path.Base(filename),
		"LogLevel":        t_LogLevel,
		"ExternalTimeout": t_ExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
	}
	return Test(opts)
}

func t_testFileOpts(filename string, t_opts map[string]interface{}) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath":   path.Dir(filename),
		"FilePath":        path.Base(filename),
		"LogLevel":        t_LogLevel,
		"ExternalTimeout": t_ExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	return Test(opts)
}

func t_testDirectory(filename string) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath":   filename,
		"LogLevel":        t_LogLevel,
		"ExternalTimeout": t_ExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
	}
	return Test(opts)
}

func t_testDirectoryOpts(filename string, t_opts map[string]interface{}) *HtmlTest {
	opts := map[string]interface{}{
		"DirectoryPath":   filename,
		"LogLevel":        t_LogLevel,
		"ExternalTimeout": t_ExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
	}
	mergo.MergeWithOverwrite(&opts, t_opts)
	return Test(opts)
}
