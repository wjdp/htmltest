package htmltest

import (
	"github.com/daviddengcn/go-assert"
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

func TestHTML5Page(t *testing.T) {
	// Page containing HTML5 tags
	hT := t_testFile("fixtures/html/html5_tags.html")
	t_expectIssueCount(t, hT, 0)
}

func TestNormalLookingPage(t *testing.T) {
	// Page containing HTML5 tags
	t_SkipShortExternal(t)
	hT := t_testFile("fixtures/html/normal_looking_page.html")
	t_expectIssueCount(t, hT, 0)
}

func TestCacheIntegration(t *testing.T) {
	t_testFileOpts("fixtures/links/linkWithHttps.html",
		map[string]interface{}{"EnableCache": true})
	hT2 := t_testFileOpts("fixtures/links/linkWithHttps.html",
		map[string]interface{}{"EnableCache": true, "NoRun": true})
	_, okY := hT2.refCache.Get("https://github.com/octocat/Spoon-Knife/issues")
	_, okN := hT2.refCache.Get("https://github.com/octocat/Spoon-Knife/milestones")
	assert.IsTrue(t, "link in cache", okY)
	assert.IsFalse(t, "link not in cache", okN)
}

func TestConcurrencyDirExternals(t *testing.T) {
	t_SkipShortExternal(t)
	hT := t_testDirectoryOpts("fixtures/concurrency/manyBrokenExt",
		map[string]interface{}{"TestFilesConcurrently": true}) // "LogLevel": 1
	t_expectIssueCount(t, hT, 26)
}
