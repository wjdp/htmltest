package htmltest

import (
	"testing"
)

// Spec tests

func TestScriptExternalSrcValid(t *testing.T) {
	// passes for valid external src
	t_testFile("fixtures/scripts/script_valid_external.html")
	t_expectIssueCount(t, 0)
}

func TestScriptExternalSrcBroken(t *testing.T) {
	// fails for broken external src
	t_testFile("fixtures/scripts/script_broken_external.html")
	t_expectIssueCount(t, 1)
	// t_expectIssue(t, "no such host", 1)
}

func TestScriptExternalInsecureDefault(t *testing.T) {
	// passes for HTTP scripts by default
	t_testFile("fixtures/scripts/scriptInsecure.html")
	t_expectIssueCount(t, 0)
}

func TestScriptExternalInsecureOption(t *testing.T) {
	// fails for HTTP scripts when asked
	t_testFileOpts("fixtures/scripts/scriptInsecure.html",
		map[string]interface{}{"EnforceHTTPS": true})
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "is not an HTTPS target", 1)
}

func TestScriptInternalSrcValid(t *testing.T) {
	// works for valid internal src
	t_testFile("fixtures/scripts/script_valid_internal.html")
	t_expectIssueCount(t, 0)
}

func TestScriptInternalSrcBroken(t *testing.T) {
	// fails for missing internal src
	t_testFile("fixtures/scripts/script_missing_internal.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "target does not exist", 1)
}

func TestScriptSrcBlank(t *testing.T) {
	// fails for blank src
	t_testFile("fixtures/scripts/scriptSrcBlank.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "src attribute present but empty", 1)
}

func TestScriptContentValid(t *testing.T) {
	// works for present content
	t_testFile("fixtures/scripts/script_content.html")
	t_expectIssueCount(t, 0)
}

func TestScriptContentAbsent(t *testing.T) {
	// fails for absent content, either content is missing or src attr missing
	t_testFile("fixtures/scripts/script_content_absent.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "script content missing / no src attribute", 1)
}

func TestScriptInPre(t *testing.T) {
	// works for broken script within pre & code
	t.Skip("TODO: ignore stuff in <pre> and <code>")
	t_testFile("fixtures/scripts/script_in_pre.html")
	t_expectIssueCount(t, 0)
}

func TestScriptSrcIgnore(t *testing.T) {
	// ignores links via url_ignore
	t.Skip("url ignoring not implemented")
	t_testFile("fixtures/scripts/ignorableLinksViaOptions.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "", 1)
}

func TestScriptIgnorable(t *testing.T) {
	t_testFile("fixtures/scripts/scriptIgnorable.html")
	t_expectIssueCount(t, 0)
}
