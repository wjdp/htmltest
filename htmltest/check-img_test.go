package htmltest

import (
	// "path"
	"testing"
)

func TestExternalImageWorking(t *testing.T) {
	// passes for existing external images
	hT := t_testFile("fixtures/images/existingImageExternal.html")
	t_expectIssueCount(t, hT, 0)
}

func TestExternalImageMissing(t *testing.T) {
	// fails for missing external images
	hT := t_testFile("fixtures/images/missingImageExternal.html")
	t_expectIssueCount(t, hT, 1)
	// Issue contains "no such host"
	// t_expectIssue(t, hT, "no such host", 1)
}

func TestExternalImageMissingProtocolValid(t *testing.T) {
	// works for valid images missing the protocol
	hT := t_testFile("fixtures/images/image_missing_protocol_valid.html")
	t_expectIssueCount(t, hT, 0)
}

func TestExternalImageMissingProtocolInvalid(t *testing.T) {
	// fails for invalid images missing the protocol
	hT := t_testFile("fixtures/images/image_missing_protocol_invalid.html")
	t_expectIssueCount(t, hT, 1)
	// t_expectIssue(t, hT, message, 1)
}

func TestExternalImageInsecureDefault(t *testing.T) {
	// passes for HTTP images by default
	hT := t_testFile("fixtures/images/src_http.html")
	t_expectIssueCount(t, hT, 0)
}

func TestExternalImageInsecureOption(t *testing.T) {
	// fails for HTTP images when asked
	hT := t_testFileOpts("fixtures/images/src_http.html",
		map[string]interface{}{"EnforceHTTPS": true})
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "is not an HTTPS target", 1)
}

func TestInternalImageAbsolute(t *testing.T) {
	// properly checks absolute images
	hT := t_testFile("fixtures/images/rootRelativeImages.html")
	t_expectIssueCount(t, hT, 0)
}

func TestInternalImageRelative(t *testing.T) {
	// properly checks relative images
	hT := t_testFile("fixtures/images/relativeToSelf.html")
	t_expectIssueCount(t, hT, 0)
}

func TestInternalImageRelativeSubfolders(t *testing.T) {
	// properly checks relative images within subfolders
	hT := t_testFile("fixtures/resources/books/nestedRelativeImages.html")
	t_expectIssueCount(t, hT, 0)
}

func TestInternalImageMissing(t *testing.T) {
	// fails for missing internal images
	hT := t_testFile("fixtures/images/missingImageInternal.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "target does not exist", 1)
}

func TestInternalImageMissingCharsAndCases(t *testing.T) {
	// fails for image with default mac filename
	hT := t_testFile("fixtures/images/terribleImageName.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "target does not exist", 1)
}

func TestInternalWithBase(t *testing.T) {
	// properly checks relative images with base
	t.Skip("base tag not supported")
	hT := t_testFile("fixtures/images/relativeWithBase.html")
	t_expectIssueCount(t, hT, 0)
}

func TestImageIgnorable(t *testing.T) {
	// ignores images marked as data-proofer-ignore
	hT := t_testFile("fixtures/images/ignorableImages.html")
	t_expectIssueCount(t, hT, 0)
}

func TestSrcMising(t *testing.T) {
	// fails for image with no src
	hT := t_testFile("fixtures/images/missingImageSrc.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "src attribute missing", 1)
}

func TestSrcEmpty(t *testing.T) {
	// fails for image with empty src
	hT := t_testFile("fixtures/images/emptyImageSrc.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "src attribute empty", 1)
}

func TestSrcLineBreaks(t *testing.T) {
	// deals with linebreaks in src
	hT := t_testFile("fixtures/images/lineBreaks.html")
	t_expectIssueCount(t, hT, 0)
}

// TODO empty src

func TestSrcIgnored(t *testing.T) {
	// ignores images via url_ignore
	t.Skip("url ignore patterns not yet implemented")
	hT := t_testFile("fixtures/images/???.html")
	t_expectIssueCount(t, hT, 0)
}

func TestSrcDataURI(t *testing.T) {
	// properly ignores data URI images
	hT := t_testFile("fixtures/images/workingDataURIImage.html")
	t_expectIssueCount(t, hT, 0)
}

func TestSrcSet(t *testing.T) {
	// works for images with a srcset
	hT := t_testFile("fixtures/images/srcSetCheck.html")
	t_expectIssueCount(t, hT, 0)
}

func TestSrcSetMissing(t *testing.T) {
	// fails for images with an alt but missing src or srcset
	hT := t_testFile("fixtures/images/srcSetMissingImage.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "src attribute missing", 1)
}

func TestSrcSetMissingAlt(t *testing.T) {
	// fails for images with a srcset but missing alt
	hT := t_testFile("fixtures/images/srcSetMissingAlt.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "alt attribute missing", 1)
}

func TestSrcSetMissingAltIgnore(t *testing.T) {
	// ignores missing alt tags when asked for srcset
	hT := t_testFileOpts("fixtures/images/srcSetIgnorable.html",
		map[string]interface{}{"IgnoreAltMissing": true})
	t_expectIssueCount(t, hT, 0)
}

func TestAltMissing(t *testing.T) {
	// fails for image without alt attribute
	hT := t_testFile("fixtures/images/missingImageAlt.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "alt attribute missing", 1)
}

func TestAltEmpty(t *testing.T) {
	// fails for image with an empty alt attribute
	hT := t_testFile("fixtures/images/missingImageAltText.html")
	t_expectIssueCount(t, hT, 1)
	t_expectIssue(t, hT, "alt text empty", 1)
}

func TestAltSpaces(t *testing.T) {
	// fails for image with nothing but spaces in alt attribute
	hT := t_testFile("fixtures/images/emptyImageAltText.html")
	t_expectIssueCount(t, hT, 3)
	t_expectIssue(t, hT, "alt text contains only whitespace", 1)
}

func TestAltIgnoreMissing(t *testing.T) {
	// ignores missing alt tags when asked
	hT := t_testFileOpts("fixtures/images/ignorableAltViaOptions.html",
		map[string]interface{}{"IgnoreAltMissing": true})
	t_expectIssueCount(t, hT, 0)
}

func TestImagePre(t *testing.T) {
	// works for broken images within pre & code
	hT := t_testFile("fixtures/images/badImagesInPre.html")
	t_expectIssueCount(t, hT, 0)
}

func TestImageMultipleProblems(t *testing.T) {
	hT := t_testFile("fixtures/images/multipleProblems.html")
	t_expectIssueCount(t, hT, 6)
	t_expectIssue(t, hT, "alt text empty", 1)
	t_expectIssue(t, hT, "target does not exist", 2)
	t_expectIssue(t, hT, "alt attribute missing", 1)
	t_expectIssue(t, hT, "src attribute missing", 1)
}
