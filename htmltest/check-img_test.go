package htmltest

import (
	// "path"
	"testing"
)

func TestImageExternalWorking(t *testing.T) {
	// passes for existing external images
	hT := tTestFileOpts("fixtures/images/existingImageExternal.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestImageExternalMissing(t *testing.T) {
	// fails for missing external images
	hT := tTestFileOpts("fixtures/images/missingImageExternal.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	// Issue contains "no such host"
	// tExpectIssue(t, hT, "no such host", 1)
}

func TestImageExternalMissingProtocolValid(t *testing.T) {
	// works for valid images missing the protocol
	hT := tTestFileOpts("fixtures/images/image_missing_protocol_valid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestImageExternalMissingProtocolInvalid(t *testing.T) {
	// fails for invalid images missing the protocol
	hT := tTestFileOpts("fixtures/images/image_missing_protocol_invalid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	// tExpectIssue(t, hT, message, 1)
}

func TestImageExternalInsecureDefault(t *testing.T) {
	// passes for HTTP images by default
	hT := tTestFileOpts("fixtures/images/src_http.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestImageExternalInsecureOption(t *testing.T) {
	// fails for HTTP images when asked
	hT := tTestFileOpts("fixtures/images/src_http.html",
		map[string]interface{}{"EnforceHTTPS": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "is not an HTTPS target", 1)
}

func TestImageInternalAbsolute(t *testing.T) {
	// properly checks absolute images
	hT := tTestFile("fixtures/images/rootRelativeImages.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageInternalRelative(t *testing.T) {
	// properly checks relative images
	hT := tTestFile("fixtures/images/relativeToSelf.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageInternalRelativeSubfolders(t *testing.T) {
	// properly checks relative images within subfolders
	hT := tTestFile("fixtures/resources/books/nestedRelativeImages.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageInternalMissing(t *testing.T) {
	// fails for missing internal images
	hT := tTestFile("fixtures/images/missingImageInternal.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestImageInternalMissingCharsAndCases(t *testing.T) {
	// fails for image with default mac filename
	hT := tTestFile("fixtures/images/terribleImageName.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestImageInternalWithBase(t *testing.T) {
	// properly checks relative images with base
	t.Skip("absolute base tags not supported")
	hT := tTestFile("fixtures/images/relativeWithBase.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageIgnorable(t *testing.T) {
	// ignores images marked as data-proofer-ignore
	hT := tTestFile("fixtures/images/ignorableImages.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageSrcMising(t *testing.T) {
	// fails for image with no src
	hT := tTestFile("fixtures/images/missingImageSrc.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "src attribute missing", 1)
}

func TestImageSrcEmpty(t *testing.T) {
	// fails for image with empty src
	hT := tTestFile("fixtures/images/emptyImageSrc.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "src attribute empty", 1)
}

func TestImageSrcLineBreaks(t *testing.T) {
	// deals with linebreaks in src
	tSkipShortExternal(t) // TODO use internal images
	hT := tTestFile("fixtures/images/lineBreaks.html")
	tExpectIssueCount(t, hT, 0)
}

// TODO empty src

func TestImageSrcIgnored(t *testing.T) {
	// ignores images via url_ignore
	t.Skip("url ignore patterns not yet implemented")
	hT := tTestFile("fixtures/images/???.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageSrcDataURI(t *testing.T) {
	// properly ignores data URI images
	hT := tTestFile("fixtures/images/workingDataURIImage.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageSrcSet(t *testing.T) {
	// works for images with a srcset
	hT := tTestFile("fixtures/images/srcSetCheck.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageSrcSetMissing(t *testing.T) {
	// fails for images with an alt but missing src or srcset
	hT := tTestFile("fixtures/images/srcSetMissingImage.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "src attribute missing", 1)
}

func TestImageSrcSetMissingAlt(t *testing.T) {
	// fails for images with a srcset but missing alt
	hT := tTestFile("fixtures/images/srcSetMissingAlt.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "alt attribute missing", 1)
}

func TestImageSrcSetMissingAltIgnore(t *testing.T) {
	// ignores missing alt tags when asked for srcset
	hT := tTestFileOpts("fixtures/images/srcSetIgnorable.html",
		map[string]interface{}{"IgnoreAltMissing": true})
	tExpectIssueCount(t, hT, 0)
}

func TestImageAltMissing(t *testing.T) {
	// fails for image without alt attribute
	hT := tTestFile("fixtures/images/missingImageAlt.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "alt attribute missing", 1)
}

func TestImageAltEmpty(t *testing.T) {
	// fails for image with an empty alt attribute
	hT := tTestFile("fixtures/images/missingImageAltText.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "alt text empty", 1)
}

func TestImageAltSpaces(t *testing.T) {
	// fails for image with nothing but spaces in alt attribute
	hT := tTestFile("fixtures/images/emptyImageAltText.html")
	tExpectIssueCount(t, hT, 3)
	tExpectIssue(t, hT, "alt text contains only whitespace", 1)
}

func TestImageAltIgnoreMissing(t *testing.T) {
	// ignores missing alt tags when asked
	hT := tTestFileOpts("fixtures/images/ignorableAltViaOptions.html",
		map[string]interface{}{"IgnoreAltMissing": true})
	tExpectIssueCount(t, hT, 0)
}

func TestImagePre(t *testing.T) {
	// works for broken images within pre & code
	hT := tTestFile("fixtures/images/badImagesInPre.html")
	tExpectIssueCount(t, hT, 0)
}

func TestImageMultipleProblems(t *testing.T) {
	hT := tTestFile("fixtures/images/multipleProblems.html")
	tExpectIssueCount(t, hT, 6)
	tExpectIssue(t, hT, "alt text empty", 1)
	tExpectIssue(t, hT, "target does not exist", 2)
	tExpectIssue(t, hT, "alt attribute missing", 1)
	tExpectIssue(t, hT, "src attribute missing", 1)
}
