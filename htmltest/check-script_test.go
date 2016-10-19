package htmltest

import (
	"testing"
)

// Spec tests

func TestScriptExternalSrcValid(t *testing.T) {
	// passes for valid external src
	hT := t_testFile("fixtures/scripts/script_valid_external.html")
	t_expectIssueCount(t, hT, 0)
}

func TestScriptExternalSrcBroken(t *testing.T) {
	// fails for broken external src
	hT := t_testFile("fixtures/scripts/script_broken_external.html")
	t_expectIssueCount(t, hT, 1)
	// t_expectIssue(t, hT, "no such host", 1)
}

func TestScriptExternalInsecureDefault(t *testing.T) {
	// passes for HTTP scripts by default
	hT := t_testFile("fixtures/scripts/scriptInsecure.html")
	t_expectIssueCount(t, hT, 0)
}

func TestScriptExternalInsecureOption(t *testing.T) {
	// fails for HTTP scripts when asked
	hT := t_testFileOpts("fixtures/scripts/scriptInsecure.html",
		map[string]interface{}{"EnforceHTTPS": true})
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "is not an HTTPS target", 1)
}

func TestScriptInternalSrcValid(t *testing.T) {
	// works for valid internal src
	hT := t_testFile("fixtures/scripts/script_valid_internal.html")
	t_expectIssueCount(t, hT, 0)
}

func TestScriptInternalSrcBroken(t *testing.T) {
	// fails for missing internal src
	hT := t_testFile("fixtures/scripts/script_missing_internal.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "target does not exist", 1)
}

func TestScriptSrcBlank(t *testing.T) {
	// fails for blank src
	hT := t_testFile("fixtures/scripts/scriptSrcBlank.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "src attribute present but empty", 1)
}

func TestScriptContentValid(t *testing.T) {
	// works for present content
	hT := t_testFile("fixtures/scripts/script_content.html")
	t_expectIssueCount(t, hT, 0)
}

func TestScriptContentAbsent(t *testing.T) {
	// fails for absent content, either content is missing or src attr missing
	hT := t_testFile("fixtures/scripts/script_content_absent.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "script content missing / no src attribute", 1)
}

func TestScriptInPre(t *testing.T) {
	// works for broken script within pre & code
	hT := t_testFile("fixtures/scripts/script_in_pre.html")
	t_expectIssueCount(t, hT, 0)
}

func TestScriptSrcIgnore(t *testing.T) {
	// ignores links via url_ignore
	t.Skip("url ignoring not implemented")
	hT := t_testFile("fixtures/scripts/ignorableLinksViaOptions.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "", 1)
}

func TestScriptIgnorable(t *testing.T) {
	hT := t_testFile("fixtures/scripts/scriptIgnorable.html")
	t_expectIssueCount(t, hT, 0)
}
