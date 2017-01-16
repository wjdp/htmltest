package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

// Kepe it quiet
const tLogLevel int = issues.LevelNone

// We're running non-concurrently, speed up the tests by turning down the
// timeout. Assumes we're on a good connection.
const tExternalTimeout int = 3

func tExpectIssue(t *testing.T, hT *HTMLTest, message string, expected int) {
	c := hT.issueStore.MessageMatchCount(message)
	if c != expected {
		hT.issueStore.DumpIssues(true)
		t.Error("expected issue", message, "count", expected, "!=", c)
	}
}

func tExpectIssueCount(t *testing.T, hT *HTMLTest, expected int) {
	c := hT.issueStore.Count(issues.LevelError)
	if c != expected {
		hT.issueStore.DumpIssues(true)
		t.Error("expected", expected, "issues,", c, "found")
	}
}

func tTestFile(filename string) *HTMLTest {
	opts := map[string]interface{}{
		"DirectoryPath":   path.Dir(filename),
		"FilePath":        path.Base(filename),
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
	return Test(opts)
}

func tTestFileOpts(filename string, tOpts map[string]interface{}) *HTMLTest {
	opts := map[string]interface{}{
		"DirectoryPath":   path.Dir(filename),
		"FilePath":        path.Base(filename),
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
	mergo.MergeWithOverwrite(&opts, tOpts)
	return Test(opts)
}

func tTestDirectory(filename string) *HTMLTest {
	opts := map[string]interface{}{
		"DirectoryPath":   filename,
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
	return Test(opts)
}

func tTestDirectoryOpts(filename string, tOpts map[string]interface{}) *HTMLTest {
	opts := map[string]interface{}{
		"DirectoryPath":   filename,
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
	mergo.MergeWithOverwrite(&opts, tOpts)
	return Test(opts)
}

func tSkipShortExternal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring network calls in short mode")
	}
}
