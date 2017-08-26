package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"testing"
)

// Keep it quiet
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

// Default options for running a file test
func defaultFileTestOpts(filename string) map[string]interface{} {
	return map[string]interface{}{
		"DirectoryPath":   path.Dir(filename),
		"FilePath":        path.Base(filename),
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
}

// Test a single file and return the run test
func tTestFile(filename string) *HTMLTest {
	return Test(defaultFileTestOpts(filename))
}

// Test a single file with custom options and return the run test
func tTestFileOpts(filename string, tOpts map[string]interface{}) *HTMLTest {
	opts := defaultFileTestOpts(filename)
	mergo.MergeWithOverwrite(&opts, tOpts)
	return Test(opts)
}

// Default options for running a directory test
func defaultDirectoryTestOpts(filename string) map[string]interface{} {
	return map[string]interface{}{
		"DirectoryPath":   filename,
		"LogLevel":        tLogLevel,
		"ExternalTimeout": tExternalTimeout,
		"EnableCache":     false,
		"EnableLog":       false,
		"CheckDoctype":    false,
	}
}

func tTestDirectory(filename string) *HTMLTest {
	return Test(defaultDirectoryTestOpts(filename))
}

func tTestDirectoryOpts(filename string, tOpts map[string]interface{}) *HTMLTest {
	opts := defaultDirectoryTestOpts(filename)
	mergo.MergeWithOverwrite(&opts, tOpts)
	return Test(opts)
}

func tSkipShortExternal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring network calls in short mode")
	}
}
