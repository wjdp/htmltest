package htmltest

import (
	"testing"
)

// Spec tests

func TestScriptExternalSrcValid(t *testing.T) {
	// passes for valid external src
	hT := tTestFileOpts("fixtures/scripts/script_valid_external.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestScriptExternalSrcBroken(t *testing.T) {
	// fails for broken external src
	hT := tTestFileOpts("fixtures/scripts/script_broken_external.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	// tExpectIssue(t, hT, "no such host", 1)
}

func TestScriptExternalInsecureDefault(t *testing.T) {
	// passes for HTTP scripts by default
	hT := tTestFileOpts("fixtures/scripts/scriptInsecure.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestScriptExternalInsecureOption(t *testing.T) {
	// fails for HTTP scripts when asked
	hT := tTestFileOpts("fixtures/scripts/scriptInsecure.html",
		map[string]interface{}{"EnforceHTTPS": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "is not an HTTPS target", 1)
}

func TestScriptInternalSrcValid(t *testing.T) {
	// works for valid internal src
	hT := tTestFile("fixtures/scripts/script_valid_internal.html")
	tExpectIssueCount(t, hT, 0)
}

func TestScriptInternalSrcBroken(t *testing.T) {
	// fails for missing internal src
	hT := tTestFile("fixtures/scripts/script_missing_internal.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestScriptSrcBlank(t *testing.T) {
	// fails for blank src
	hT := tTestFile("fixtures/scripts/scriptSrcBlank.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "src attribute present but empty", 1)
}

func TestScriptContentValid(t *testing.T) {
	// works for present content
	hT := tTestFile("fixtures/scripts/script_content.html")
	tExpectIssueCount(t, hT, 0)
}

func TestScriptContentAbsent(t *testing.T) {
	// fails for absent content, either content is missing or src attr missing
	hT := tTestFile("fixtures/scripts/script_content_absent.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "script content missing / no src attribute", 1)
}

func TestScriptInPre(t *testing.T) {
	// works for broken script within pre & code
	hT := tTestFile("fixtures/scripts/script_in_pre.html")
	tExpectIssueCount(t, hT, 0)
}

func TestScriptSrcIgnore(t *testing.T) {
	// ignores links via url_ignore
	t.Skip("url ignoring not implemented")
	hT := tTestFile("fixtures/scripts/ignorableLinksViaOptions.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "", 1)
}

func TestScriptIgnorable(t *testing.T) {
	hT := tTestFile("fixtures/scripts/scriptIgnorable.html")
	tExpectIssueCount(t, hT, 0)
}

func TestScriptIgnorableChildren(t *testing.T) {
	hT := tTestFile("fixtures/scripts/scriptIgnorableChildren.html")
	tExpectIssueCount(t, hT, 0)
}
