package htmltest

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestMissingOptions(t *testing.T) {
	// returns error when we don't set FilePath nor DirectoryPath
	_, err := Test(map[string]interface{}{})
	assert.NotEquals(t, "Error", err, nil)
	assert.Equals(t, "Error", err.Error(),
		"Neither FilePath nor DirectoryPath provided")
}

func TestDirectoryPathNonExistent(t *testing.T) {
	// returns error when DirectoryPath doesn't actually exist
	_, err := Test(map[string]interface{}{
		"DirectoryPath": "no-exist",
	})
	assert.NotEquals(t, "Error", err, nil)
	assert.Equals(t, "Error", err.Error(),
		"Cannot access 'no-exist', no such directory.")
}

func TestDirectoryPathIsFile(t *testing.T) {
	// returns error DirectoryPath resolves to a file
	_, err := Test(map[string]interface{}{
		"DirectoryPath": "fixtures/utils/file",
	})
	assert.NotEquals(t, "Error", err, nil)
	assert.Equals(t, "Error", err.Error(),
		"DirectoryPath 'fixtures/utils/file' is a file, not a directory.")
}

func TestFilePathMissing(t *testing.T) {
	// returns error when we can't find FilePath
	_, err := Test(map[string]interface{}{
		"DirectoryPath": "fixtures/utils",
		"FilePath":      "no-file",
	})
	assert.NotEquals(t, "Error", err, nil)
	assert.Equals(t, "Error", err.Error(),
		"Could not find FilePath 'no-file' in 'fixtures/utils'")
}

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

func TestCountDocuments(t *testing.T) {
	hT := tTestDirectory("fixtures/documents/folder-ok")
	assert.Equals(t, "CountDocuments", hT.CountDocuments(), 3)
}

func TestCountErrors(t *testing.T) {
	hT := tTestDirectory("fixtures/documents/folder-not-ok")
	assert.Equals(t, "CountErrors", hT.CountErrors(), 2)
}

func TestFileExtensionDefault(t *testing.T) {
	// Non .html files are ignored
	hT := tTestDirectory("fixtures/documents/folder-htm")
	assert.Equals(t, "CountDocuments", hT.CountDocuments(), 0)
	assert.Equals(t, "CountErrors", hT.CountErrors(), 0)
}

func TestFileExtensionOption(t *testing.T) {
	// FileExtension (+DirectoryIndex) works when set
	hT := tTestDirectoryOpts("fixtures/documents/folder-htm", map[string]interface{}{
		"FileExtension":  ".htm",
		"DirectoryIndex": "index.htm",
	})
	assert.Equals(t, "CountDocuments", hT.CountDocuments(), 3)
	tExpectIssueCount(t, hT, 1)
}

func TestCacheIntegration(t *testing.T) {
	tSkipShortExternal(t)
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

func TestRedirectLimitDefault(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/http_no_redirect.html",
		map[string]interface{}{"RedirectLimit": -2})
	tExpectIssueCount(t, hT, 0)
	hT = tTestFileOpts("fixtures/links/http_one_redirect.html",
		map[string]interface{}{"RedirectLimit": -1})
	tExpectIssueCount(t, hT, 0)
}

func TestRedirectLimitOk(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/http_no_redirect.html",
		map[string]interface{}{"RedirectLimit": 0})
	tExpectIssueCount(t, hT, 0)
	hT = tTestFileOpts("fixtures/links/http_one_redirect.html",
		map[string]interface{}{"RedirectLimit": 1})
	tExpectIssueCount(t, hT, 0)
}

func TestRedirectLimitExceeded(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/http_one_redirect.html",
		map[string]interface{}{"RedirectLimit": 0})
	tExpectIssueCount(t, hT, 1)
}
