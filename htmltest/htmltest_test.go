package htmltest

import (
	"testing"
)

func TestCheckAnchorsDisable(t *testing.T) {
	hT := t_testFileOpts("fixtures/links/brokenLinkInternal.html",
		map[string]interface{}{"CheckAnchors": false})
	t_expectIssueCount(t, hT, 0)
}

func TestCheckLinksDisable(t *testing.T) {
	hT := t_testFileOpts("fixtures/links/head_link_href_absent.html",
		map[string]interface{}{"CheckLinks": false})
	t_expectIssueCount(t, hT, 0)
}

func TestCheckImagesDisable(t *testing.T) {
	hT := t_testFileOpts("fixtures/images/emptyImageSrc.html",
		map[string]interface{}{"CheckImages": false})
	t_expectIssueCount(t, hT, 0)
}

func TestCheckScriptsDisable(t *testing.T) {
	hT := t_testFileOpts("fixtures/scripts/script_content_absent.html",
		map[string]interface{}{"CheckScripts": false})
	t_expectIssueCount(t, hT, 0)
}
