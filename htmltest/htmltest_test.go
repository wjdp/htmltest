package htmltest

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestCheckAnchorsDisable(t *testing.T) {
	hT := tTestFileOpts("fixtures/links/brokenLinkInternal.html",
		map[string]interface{}{"CheckAnchors": false})
	tExpectIssueCount(t, hT, 0)
}

func TestCheckLinksDisable(t *testing.T) {
	hT := tTestFileOpts("fixtures/links/head_link_href_absent.html",
		map[string]interface{}{"CheckLinks": false})
	tExpectIssueCount(t, hT, 0)
}

func TestCheckImagesDisable(t *testing.T) {
	hT := tTestFileOpts("fixtures/images/emptyImageSrc.html",
		map[string]interface{}{"CheckImages": false})
	tExpectIssueCount(t, hT, 0)
}

func TestCheckScriptsDisable(t *testing.T) {
	hT := tTestFileOpts("fixtures/scripts/script_content_absent.html",
		map[string]interface{}{"CheckScripts": false})
	tExpectIssueCount(t, hT, 0)
}

func TestHTML5Page(t *testing.T) {
	// Page containing HTML5 tags
	hT := tTestFile("fixtures/html/html5_tags.html")
	tExpectIssueCount(t, hT, 0)
}

func TestNormalLookingPage(t *testing.T) {
	// Page containing HTML5 tags
	hT := tTestFileOpts("fixtures/html/normal_looking_page.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestCacheIntegration(t *testing.T) {
	tTestFileOpts("fixtures/links/https-valid.html",
		map[string]interface{}{"EnableCache": true})
	hT2 := tTestFileOpts("fixtures/links/https-valid.html",
		map[string]interface{}{"EnableCache": true, "NoRun": true})
	_, okY := hT2.refCache.Get("https://github.com/octocat/Spoon-Knife/issues")
	_, okN := hT2.refCache.Get("https://github.com/octocat/Spoon-Knife/milestones")
	assert.IsTrue(t, "link in cache", okY)
	assert.IsFalse(t, "link not in cache", okN)
}

func TestConcurrencyDirExternals(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestDirectoryOpts("fixtures/concurrency/manyBrokenExt",
		map[string]interface{}{"TestFilesConcurrently": true}) // "LogLevel": 1
	tExpectIssueCount(t, hT, 26)
}
