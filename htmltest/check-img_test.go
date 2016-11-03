package htmltest

import (
	// "path"
	"testing"
)

func TestExternalImageWorking(t *testing.T) {
	// passes for existing external images
	t_testFile("fixtures/images/existingImageExternal.html")
	t_expectIssueCount(t, 0)
}

func TestExternalImageMissing(t *testing.T) {
	// fails for missing external images
	t_testFile("fixtures/images/missingImageExternal.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestExternalImageMissingProtocolValid(t *testing.T) {
	// works for valid images missing the protocol
	t_testFile("fixtures/images/image_missing_protocol_valid.html")
	t_expectIssueCount(t, 0)
}

func TestExternalImageMissingProtocolInvalid(t *testing.T) {
	// fails for invalid images missing the protocol
	t_testFile("fixtures/images/image_missing_protocol_invalid.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, message, 1)
}

func TestExternalImageInsecureDefault(t *testing.T) {
	// passes for HTTP images by default
	t_testFile("fixtures/images/src_http.html")
	t_expectIssueCount(t, 0)
}

func TestExternalImageInsecureOption(t *testing.T) {
	// fails for HTTP images when asked
	t_testFile("fixtures/images/src_http.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestInternalImageAbsolute(t *testing.T) {
	// properly checks absolute images
	t_testFile("fixtures/images/rootRelativeImages.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestInternalImageRelative(t *testing.T) {
	// properly checks relative images
	t_testFile("fixtures/images/relativeToSelf.html")
	t_expectIssueCount(t, 0)
}

func TestInternalImageRelativeSubfolders(t *testing.T) {
	// properly checks relative images within subfolders
	t_testFile("fixtures/resources/books/nestedRelativeImages.html")
	t_expectIssueCount(t, 0)
}

func TestInternalImageMissing(t *testing.T) {
	// fails for missing internal images
	t_testFile("fixtures/images/missingImageInternal.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestInternalImageMissingCharsAndCases(t *testing.T) {
	// fails for image with default mac filename
	t_testFile("fixtures/images/terribleImageName.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestInternalWithBase(t *testing.T) {
	// properly checks relative images with base
	t.Skip("base tag not supported")
	t_testFile("fixtures/images/relativeWithBase.html")
}

func TestSrcMising(t *testing.T) {
	// fails for image with no src
	t_testFile("fixtures/images/missingImageSrc.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestSrcIgnored(t *testing.T) {
	// ignores images via url_ignore
	t.Skip("url ignore patterns not yet implemented")
	t_testFile("fixtures/images/ignorableImages.html")
	t_expectIssueCount(t, 0)
}

func TestSrcDataURI(t *testing.T) {
	// properly ignores data URI images
	t_testFile("fixtures/images/workingDataURIImage.html")
	t_expectIssueCount(t, 0)
}

func TestSrcSet(t *testing.T) {
	// works for images with a srcset
	t_testFile("fixtures/images/srcSetCheck.html")
	t_expectIssueCount(t, 0)
}

func TestSrcSetMissing(t *testing.T) {
	// fails for images with an alt but missing src or srcset
	t_testFile("fixtures/images/srcSetMissingImage.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, message, 1)
}

func TestSrcSetMissingAlt(t *testing.T) {
	// fails for images with a srcset but missing alt
	t_testFile("fixtures/images/srcSetMissingAlt.html")
	t_expectIssueCount(t, 1)
	t_expectIssue(t, "missing alt", 1)
}

func TestSrcSetMissingAltIgnore(t *testing.T) {
	// ignores missing alt tags when asked for srcset
	t.Skip("New option needed")
	t_testFile("fixtures/images/srcSetIgnorable.html")
	t_expectIssueCount(t, 0)
}

func TestAltTextMissing(t *testing.T) {
	// fails for image without alt attribute
	t_testFile("fixtures/images/missingImageAlt.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestAltTextEmpty(t *testing.T) {
	// fails for image with an empty alt attribute
	t_testFile("fixtures/images/missingImageAltText.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestAltTextSpaces(t *testing.T) {
	// fails for image with nothing but spaces in alt attribute
	t_testFile("fixtures/images/emptyImageAltText.html")
	t_expectIssueCount(t, 0)
	t_expectIssue(t, message, 1)
}

func TestAltTextIgnoreMissing(t *testing.T) {
	// ignores missing alt tags when asked
	t.Skip("New option needed")
	t_testFile("fixtures/images/ignorableAltViaOptions.html")
	t_expectIssueCount(t, 0)
}

func TestAltTextIgnoreEmpty(t *testing.T) {
	// ignores missing alt attribute when asked
	t.Skip("New option needed")
	t_testFile("fixtures/images/missingImageAlt.html")
	t_expectIssueCount(t, 0)
}
