package htmltest

import (
	"testing"
)

// Passes for valid doctype
func TestDoctypeValid(t *testing.T) {
	hT1 := tTestFileOpts("fixtures/doctype/doctype-html5.html",
		map[string]interface{}{"CheckDoctype": true})
	tExpectIssueCount(t, hT1, 0)
	hT2 := tTestFileOpts("fixtures/doctype/doctype-html4.html",
		map[string]interface{}{"CheckDoctype": true})
	tExpectIssueCount(t, hT2, 0)
	hT3 := tTestFileOpts("fixtures/doctype/doctype-xhtml.html",
		map[string]interface{}{"CheckDoctype": true})
	tExpectIssueCount(t, hT3, 0)
}

// Fails for missing doctype
func TestDoctypeMissing(t *testing.T) {
	hT := tTestFileOpts("fixtures/doctype/doctype-missing.html",
		map[string]interface{}{"CheckDoctype": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "missing doctype", 1)
}

// Fails for a doctype mixed in with page
func TestDoctypeNotFirst(t *testing.T) {
	hT := tTestFileOpts("fixtures/doctype/doctype-not-first.html",
		map[string]interface{}{"CheckDoctype": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "missing doctype", 1)
}

// Passes for html5 doctype when asked
func TestDoctypeHTML5Valid(t *testing.T) {
	hT := tTestFileOpts("fixtures/doctype/doctype-html5.html",
		map[string]interface{}{"CheckDoctype": true, "EnforceHTML5": true})
	tExpectIssueCount(t, hT, 0)
}

// Fails for other doctype when asked for html5
func TestDoctypeHTML5Invalid(t *testing.T) {
	hT := tTestFileOpts("fixtures/doctype/doctype-html4.html",
		map[string]interface{}{"CheckDoctype": true, "EnforceHTML5": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "doctype isn't html5", 1)
}
