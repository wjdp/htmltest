package htmltest

import (
	"testing"
)

// Passes for valid meta refresh without url.
func TestMetaRefreshSelfValid(t *testing.T) {
	hT := tTestFile("fixtures/meta/refresh-refresh.html")
	tExpectIssueCount(t, hT, 0)
}

// Passes for valid external meta refresh.
func TestMetaRefreshExternalValid(t *testing.T) {
	hT := tTestFileOpts("fixtures/meta/refresh-external-valid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

// Passes for valid external meta refresh with space between the ; and the
// url def.
func TestMetaRefreshExternalValidCaps(t *testing.T) {
	hT := tTestFileOpts("fixtures/meta/refresh-external-valid-space.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

// Passes for valid external meta refresh with URL in caps.
func TestMetaRefreshExternalValidSpace(t *testing.T) {
	hT := tTestFileOpts("fixtures/meta/refresh-external-valid-caps.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

// Fails broken external URL in meta refresh.
func TestMetaRefreshExternalBroken(t *testing.T) {
	hT := tTestFileOpts("fixtures/meta/refresh-external-broken.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
}

// Passes for valid internal meta refresh.
func TestMetaRefreshInternalValid(t *testing.T) {
	hT := tTestFile("fixtures/meta/refresh-internal-valid.html")
	tExpectIssueCount(t, hT, 0)
}

// Fails for broken internal path in meta refresh.
func TestMetaRefreshInternalBroken(t *testing.T) {
	hT := tTestFile("fixtures/meta/refresh-internal-broken.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

// Fails for missing content attribute when http-equiv="refresh" present.
func TestMetaRefreshContentMissing(t *testing.T) {
	hT := tTestFile("fixtures/meta/refresh-content-missing.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "missing content attribute in meta refresh", 1)
}

// Fails for blank content attribute when http-equiv="refresh" present.
func TestMetaRefreshContentBlank(t *testing.T) {
	hT := tTestFile("fixtures/meta/refresh-content-blank.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "blank content attribute in meta refresh", 1)
}

// Fails for invalid content attribute when http-equiv="refresh" present.
// The attribute should be a positive integer and may be suffixed with ;url=
// and a path.
func TestMetaRefreshContentInvalid(t *testing.T) {
	// Invalid when straight refresh
	hT1 := tTestFile("fixtures/meta/refresh-content-invalid-refresh.html")
	tExpectIssueCount(t, hT1, 1)
	tExpectIssue(t, hT1, "invalid content attribute in meta refresh", 1)
	// Invalid when a redirect
	hT2 := tTestFile("fixtures/meta/refresh-content-invalid-redirect.html")
	tExpectIssueCount(t, hT2, 1)
	tExpectIssue(t, hT2, "invalid content attribute in meta refresh", 1)
	// Malformed separator
	hT3 := tTestFile("fixtures/meta/refresh-content-invalid-redirect.html")
	tExpectIssueCount(t, hT3, 1)
	tExpectIssue(t, hT3, "invalid content attribute in meta refresh", 1)
}

// Fails when meta refresh url starts with a single or double quote
func TestIssue92(t *testing.T) {
	hT := tTestFile("fixtures/meta/issues/92.html")
	tExpectIssueCount(t, hT, 2)
	tExpectIssue(t, hT, "url in meta refresh must not start with single or double quote", 3)
}
