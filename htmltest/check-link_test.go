package htmltest

import (
	"testing"

	"github.com/wjdp/htmltest/issues"
)

// Spec tests

func TestAnchorMissingHref(t *testing.T) {
	// fails for link with no href
	hT := tTestFile("fixtures/links/missingLinkHref.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "<a> href blank body=\"No link, man\"", 1)
}

func TestAnchorIgnorable(t *testing.T) {
	// ignores links marked as ignore data-proofer-ignore
	hT := tTestFile("fixtures/links/ignorableLinks.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorIgnorableChildren(t *testing.T) {
	// ignores links marked as ignore data-proofer-ignore in child elements
	hT := tTestFile("fixtures/links/ignorableLinksChildren.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorMatchIgnore(t *testing.T) {
	// ignores links in IgnoreURLs
	hT := tTestFileOpts("fixtures/links/brokenLinkExternalSingle.html",
		map[string]interface{}{
			"IgnoreURLs": []interface{}{"www.asdo3IRJ395295jsingrkrg4.com"},
		})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalBroken(t *testing.T) {
	// fails for broken external links
	hT := tTestFileOpts("fixtures/links/brokenLinkExternal.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
}

func TestAnchorExternalBrokenOption(t *testing.T) {
	// passes for broken external links when asked
	hT := tTestFileOpts("fixtures/links/brokenLinkExternal.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalBrokenNoVCR(t *testing.T) {
	// fails for broken external links without VCR. This is needed as the code that handles 'dial tcp' errors doesn't
	// get called with VCR. It returns a rather empty response with status code of 0.
	tSkipShortExternal(t)
	hT := tTestFile("fixtures/links/brokenLinkExternal.html")
	tExpectIssueCount(t, hT, 1)
}

func TestAnchorExternalBrokenOptionNoVCR(t *testing.T) {
	// passes for broken external links without VCR when asked.
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/brokenLinkExternal.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalIgnore(t *testing.T) {
	// ignores external links when asked
	hT := tTestFileOpts("fixtures/links/brokenLinkExternal.html",
		map[string]interface{}{"CheckExternal": false, "VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHashBrokenDefault(t *testing.T) {
	// passes for broken external hashes by default
	hT := tTestFileOpts("fixtures/links/brokenHashOnTheWeb.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHashBrokenOption(t *testing.T) {
	// fails for broken external hashes when asked
	t.Skip("Not yet implemented")
	hT := tTestFileOpts("fixtures/links/brokenHashOnTheWeb.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "no such hash", 1)
}

func TestAnchorExternalCache(t *testing.T) {
	// does not check links with parameters multiple times
	// TODO check cache is being checked
	t.Skip("Not yet implemented")
	hT := tTestFile("fixtures/links/check_just_once.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHrefMalformed(t *testing.T) {
	// does not explode on bad external links in files
	hT := tTestFileOpts("fixtures/links/bad_external_links.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 2)
}

func TestAnchorExternalInsecureDefault(t *testing.T) {
	// passes for non-HTTPS links when not asked
	hT := tTestFileOpts("fixtures/links/non_https.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalInsecureOption(t *testing.T) {
	// fails for non-HTTPS links when asked
	hT := tTestFileOpts("fixtures/links/non_https.html",
		map[string]interface{}{"EnforceHTTPS": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "is not an HTTPS target", 1)
}

func TestAnchorExternalBrokenOptionInsecure(t *testing.T) {
	// passes for non-HTTPS links when asked
	hT := tTestFileOpts("fixtures/links/non_https.html",
		map[string]interface{}{"EnforceHTTPS": true, "IgnoreExternalBrokenLinks": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalInsecureOptionIgnored(t *testing.T) {
	// passes when checking for non-HTTPS links but they're in the IgnoreURLs list
	hT := tTestFileOpts("fixtures/links/issues/94.html",
		map[string]interface{}{
			"EnforceHTTPS": true,
			"IgnoreURLs":   []interface{}{"plantuml.com", "plantuml.net", "forum.plantuml.net"},
		})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalInsecureOptionIgnoredInsecure(t *testing.T) {
	// checks non-HTTPS links when they're in the IgnoreHTTPS list
	hT := tTestFileOpts("fixtures/links/non_https_ignore.html",
		map[string]interface{}{
			"EnforceHTTPS": true,
			"IgnoreHTTPS":  []interface{}{`ben\.balter\.com`, `doesntexist\.io`},
			"VCREnable":    true,
		})
	tExpectIssueCount(t, hT, 2)
	tExpectIssue(t, hT, "is not an HTTPS target", 1)
	tExpectIssue(t, hT, "no such host", 1)
}

func TestAnchorExternalHrefIP(t *testing.T) {
	// fails for broken IP address links
	hT := tTestFileOpts("fixtures/links/ip_href.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 2)
}

func TestAnchorExternalBrokenOptionHrefIP(t *testing.T) {
	// passes for broken IP address links when asked
	hT := tTestFileOpts("fixtures/links/ip_href.html",
		map[string]interface{}{"VCREnable": true, "IgnoreExternalBrokenLinks": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHrefIPTimeout(t *testing.T) {
	// fails for broken IP address links
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/ip_timeout.html",
		map[string]interface{}{"ExternalTimeout": 1})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "request exceeded our ExternalTimeout", 1)
}

func TestAnchorExternalBrokenOptionHrefIPTimeout(t *testing.T) {
	// passes for broken IP address links when aksed
	hT := tTestFileOpts("fixtures/links/ip_timeout.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true, "ExternalTimeout": 1})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalFollowRedirects(t *testing.T) {
	// should follow redirects
	hT := tTestFileOpts("fixtures/links/linkWithRedirect.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalFollowRedirectsDisabled(t *testing.T) {
	// fails on redirects if not following
	t.Skip("Not yet implemented, need new link, times out")
	hT := tTestFileOpts("fixtures/links/linkWithRedirect.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 99)
	tExpectIssue(t, hT, "PLACEHOLDER", 99)
}

func TestAnchorExternalHTTPS(t *testing.T) {
	// should understand https
	hT := tTestFileOpts("fixtures/links/https-valid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHTTPSInvalid(t *testing.T) {
	// should understand https
	hT := tTestFileOpts("fixtures/links/https-invalid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 6)
}

func TestAnchorExternalBrokenOptionHTTPSInvalid(t *testing.T) {
	// should pass for invalid https when asked
	hT := tTestFileOpts("fixtures/links/https-invalid.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHTTPSMissingChain(t *testing.T) {
	// TODO: remove skip when incomplete-chain.badssl.com has an in-date cert
	t.Skip("The cert on incomplete-chain.badssl.com has expired, cannot use for testing at present")
	// should support https aia
	// see issue #130
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/https-incomplete-chain.html",
		map[string]interface{}{"VCREnable": false})
	tExpectIssue(t, hT, "incomplete certificate chain", 1)
}

func TestAnchorExternalBrokenOptionHTTPSMissingChain(t *testing.T) {
	// should pass for incomplete chains when asked
	hT := tTestFileOpts("fixtures/links/https-incomplete-chain.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true, "VCREnable": false})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHTTPSBadH2(t *testing.T) {
	// should connect to servers with bad http/2 support
	// See issue #49
	hT := tTestFileOpts("fixtures/links/https-valid-h2.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalRequiresAccepts(t *testing.T) {
	// should connect to servers with bad http/2 support
	// See issue #49
	hT := tTestFileOpts("fixtures/links/http_requires_accept.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalMissingProtocolValid(t *testing.T) {
	// works for valid links missing the protocol
	hT := tTestFileOpts("fixtures/links/link_missing_protocol_valid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalMissingProtocolInvalid(t *testing.T) {
	// fails for invalid links missing the protocol
	hT := tTestFileOpts("fixtures/links/link_missing_protocol_invalid.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 1)
	// tExpectIssue(t, hT, "no such host", 1)
}

func TestAnchorExternalBrokenOptionMissingProtocol(t *testing.T) {
	// passes for invalid links missing the protocol when asked
	hT := tTestFileOpts("fixtures/links/link_missing_protocol_invalid.html",
		map[string]interface{}{"IgnoreExternalBrokenLinks": true, "VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestLinkExternalHrefPipes(t *testing.T) {
	// works for pipes in the URL
	hT := tTestFileOpts("fixtures/links/escape_pipes.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHrefNonstandardChars(t *testing.T) {
	// passes non-standard characters
	hT := tTestFileOpts("fixtures/links/non_standard_characters.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalHrefUTF8(t *testing.T) {
	// passes for external UTF-8 links
	hT := tTestFileOpts("fixtures/links/utf8Link.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalInvalidBrackets(t *testing.T) {
	// passes for external UTF-8 links
	hT := tTestFile("fixtures/links/href_brackets.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "bad reference", 1)
}

func TestAnchorExternalQueryStringDefault(t *testing.T) {
	// passes when ignoring from default list of query string exempt URLs
	tSkipShortExternal(t)
	hT := tTestFile("fixtures/links/query_strings.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorExternalQueryStripQueryExcludesEmpty(t *testing.T) {
	// fails when StripQueryExcludes blank and URL doesn't like query string hits
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/query_strings.html",
		map[string]interface{}{"StripQueryExcludes": []interface{}{}})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "Non-OK status: 400", 1)
}

func TestAnchorExternalQueryStringStripQueryExcludesDiffers(t *testing.T) {
	// fails when StripQueryExcludes does not include URL and URL doesn't like query string hits
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/query_strings.html",
		map[string]interface{}{"StripQueryExcludes": []interface{}{"example.com", "test.invalid"}})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "Non-OK status: 400", 1)
}

func TestAnchorInternalBroken(t *testing.T) {
	// fails for broken internal links
	hT := tTestFile("fixtures/links/brokenLinkInternal.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestAnchorInternalBrokenIgnore(t *testing.T) {
	// fails for broken internal links
	hT := tTestFileOpts("fixtures/links/brokenLinkInternal.html",
		map[string]interface{}{"CheckInternal": false})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalIgnoreUrl(t *testing.T) {
	// ignores internal links in IgnoreURLs
	hT := tTestFileOpts("fixtures/links/brokenLinkInternal.html",
		map[string]interface{}{
			"IgnoreURLs": []interface{}{"no\\w+.html"},
		})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorFileIgnoreUrl(t *testing.T) {
	// ignores file links in IgnoreURLs
	hT := tTestFileOpts("fixtures/links/brokenLinkFile.html",
		map[string]interface{}{
			"IgnoreURLs": []interface{}{"no\\w+.html"},
		})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalRelativeLinksBase(t *testing.T) {
	// passes for relative links with a base
	hT := tTestFile("fixtures/links/relativeLinksWithBase.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorHashInternalValid(t *testing.T) {
	// passes for valid internal hash
	hT := tTestFile("fixtures/links/hashInternalOk.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorHashInternalBroken(t *testing.T) {
	// fails for broken internal hash
	hT := tTestFile("fixtures/links/hashInternalBroken.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "hash does not exist", 1)
}

func TestAnchorHashSelfValid(t *testing.T) {
	// passes for valid self hash
	hT := tTestFile("fixtures/links/hashSelfOk.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorHashSelfBroken(t *testing.T) {
	// fails for broken self hash
	hT := tTestFile("fixtures/links/hashSelfBroken.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "hash does not exist", 1)
}

func TestAnchorHashBrokenIgnore(t *testing.T) {
	// fails for broken internal hash
	hT1 := tTestFileOpts("fixtures/links/hashInternalBroken.html",
		map[string]interface{}{"CheckInternalHash": false})
	hT2 := tTestFileOpts("fixtures/links/hashSelfBroken.html",
		map[string]interface{}{"CheckInternalHash": false})
	tExpectIssueCount(t, hT1, 0)
	tExpectIssueCount(t, hT2, 0)
}

func TestAnchorHashInPre163(t *testing.T) {
	// passes for valid self hash
	hT := tTestFile("fixtures/links/hash-in-pre-163.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryRootResolve(t *testing.T) {
	// properly resolves implicit /index.html in link paths
	hT := tTestFile("fixtures/links/linkToFolder.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryRootResolveWithIgnoredDir(t *testing.T) {
	// ignoring the target of the link does not break
	hT := tTestFileOpts("fixtures/links/linkToFolder.html",
		map[string]interface{}{"IgnoreDirs": []interface{}{"folder"}})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryCustomRoot(t *testing.T) {
	// works for custom directory index file
	t.Skip("Not yet implemented")
	hT := tTestFile("fixtures/links/link_pointing_to_directory.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryCustomRootBroken(t *testing.T) {
	// fails if custom directory index file doesn't exist
	hT := tTestFile("fixtures/links/link_pointing_to_directory.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target is a directory, no index", 1)
}

func TestAnchorDirectoryNoTrailingSlash(t *testing.T) {
	// fails for internal linking to a directory without trailing slash
	hT := tTestFile("fixtures/links/link_directory_without_slash.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target is a directory, href lacks trailing slash", 1)
}

func TestAnchorDirectoryNoTrailingSlashOption(t *testing.T) {
	// passes for internal linking to a directory without trailing slash when asked
	hT := tTestFileOpts("fixtures/links/link_directory_without_slash.html",
		map[string]interface{}{"IgnoreDirectoryMissingTrailingSlash": true})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryQueryHash(t *testing.T) {
	// passes for internal linking to a directory with trailing slash
	hT := tTestFile("fixtures/links/link_directory_with_slash_query_hash.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryRootExplicit(t *testing.T) {
	// passes when linking explicitly to the directory root
	hT := tTestFile("fixtures/links/root-link-explicit.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryHtmlExtension(t *testing.T) {
	// works for custom directory index file
	hT := tTestDirectory("fixtures/links/_site/")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorDirectoryWithEncodedCharacters(t *testing.T) {
	// passes for folder with encoded characters
	hT := tTestFile("fixtures/links/linkToFolderWithSpace.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalRootLink(t *testing.T) {
	// properly checks links to root
	hT := tTestFile("fixtures/links/rootLink/rootLink.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalRelativeLinks(t *testing.T) {
	// properly checks relative links
	hT := tTestFile("fixtures/links/relativeLinks.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalHrefUrlEncoded(t *testing.T) {
	// passes for urlencoded href
	hT := tTestFile("fixtures/links/urlencoded-href.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorErrorDuplication(t *testing.T) {
	// does not dupe errors
	hT := tTestFile("fixtures/links/nodupe.html")
	tExpectIssueCount(t, hT, 1)
}

func TestAnchorInternalDashedAttrs(t *testing.T) {
	// does not complain for files with attributes containing dashes
	hT := tTestFile("fixtures/links/attributeWithDash.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalCaseMismatch(t *testing.T) {
	// does not complain for internal hash links with mismatched cases
	t.Skip("Unsure on whether we should ignore case, pretty sure we shouldn't")
	hT := tTestFile("fixtures/links/ignores_cases.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalHashBlankDefault(t *testing.T) {
	// fails for href="#" when not asked
	hT := tTestFile("fixtures/links/hash_href.html")
	tExpectIssue(t, hT, "<a> empty hash", 1)
	tExpectIssueCount(t, hT, 1)
}

func TestAnchorInternalHashBlankOption(t *testing.T) {
	// passes for href="#" when asked, see
	// https://github.com/wjdp/htmltest/issues/30
	hT1 := tTestFileOpts("fixtures/links/hash_href.html",
		map[string]interface{}{"CheckInternalHash": false})
	tExpectIssueCount(t, hT1, 0)
	hT2 := tTestFileOpts("fixtures/links/hash_href.html",
		map[string]interface{}{"IgnoreInternalEmptyHash": true})
	tExpectIssueCount(t, hT2, 0)
}

func TestAnchorInternalHashWeird(t *testing.T) {
	// works for internal links to weird encoding IDs
	hT := tTestFile("fixtures/links/encodingLink.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalUrlDoesNotExist(t *testing.T) {
	// fails for internal link not in IgnoreInternalURLs (#168)
	hT := tTestFile("fixtures/links/link_directory_internal_invalid.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestAnchorInternalUrlIgnoreUsingIgnoreInternalURLs(t *testing.T) {
	// passes for internal link in IgnoreInternalURLs option (#168)
	hT := tTestFileOpts("fixtures/links/link_directory_internal_valid.html",
		map[string]interface{}{"IgnoreInternalURLs": []interface{}{"/misc/js/script.js"}})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalUrlIgnoreInternalURLsIsStrict(t *testing.T) {
	// fails as IgnoreInternalURLs requires a string match
	hT := tTestFileOpts("fixtures/links/link_directory_internal_valid.html",
		map[string]interface{}{"IgnoreInternalURLs": []interface{}{"misc/js/script.js"}})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestAnchorInternalUrlIgnoreUsingIgnoreURLs(t *testing.T) {
	// passes for internal link in IgnoreURLs option using regex
	hT := tTestFileOpts("fixtures/links/link_directory_internal_valid.html",
		map[string]interface{}{"IgnoreURLs": []interface{}{"^/misc/js/script.js$"}})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalUrlNotIgnoreUsingIgnoreURLs(t *testing.T) {
	// fails for internal link in when doesn't match strict regex in IgnoreURLs
	hT := tTestFileOpts("fixtures/links/link_directory_internal_invalid.html",
		map[string]interface{}{"IgnoreURLs": []interface{}{"^/misc/js/script.js$"}})
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestAnchorMultipleProblems(t *testing.T) {
	// finds a mix of broken and unbroken links
	t.Skip("Only single problem, and an hash which is not yet supported.")
	// TODO make our own multiple problem file
	hT := tTestFile("fixtures/links/multipleProblems.html")
	tExpectIssueCount(t, hT, 99)
	tExpectIssue(t, hT, "PLACEHOLDER", 99)
}

func TestAnchorJavascriptLinkIgnore(t *testing.T) {
	// ignores javascript links
	hT := tTestFile("fixtures/links/javascript_link.html")
	tExpectIssueCount(t, hT, 0)
}

func TestMailtoValid(t *testing.T) {
	// ignores valid mailto links
	hT := tTestFile("fixtures/links/mailto_link.html")
	tExpectIssueCount(t, hT, 0)
}

func TestMailtoEncoded(t *testing.T) {
	// ignores valid mailto links
	hT := tTestFile("fixtures/links/mailto_encoded.html")
	tExpectIssueCount(t, hT, 0)
}

func TestMailtoEncodedInvalid(t *testing.T) {
	// ignores valid mailto links
	hT := tTestFile("fixtures/links/mailto_encoded_invalid.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "cannot decode email (invalid URL escape \"%ZZ\")", 1)
}

func TestMailtoBlank(t *testing.T) {
	// fails for blank mailto links
	hT := tTestFile("fixtures/links/blank_mailto_link.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "mailto is empty", 1)
}

func TestMailtoInvalid(t *testing.T) {
	// fails for invalid mailto links
	hT := tTestFile("fixtures/links/invalid_mailto_link.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "invalid email address (invalid format): 'octocat'", 1)
}

func TestMailtoIgnore(t *testing.T) {
	// ignores mailto links when told to
	hT := tTestFileOpts("fixtures/links/blank_mailto_link.html",
		map[string]interface{}{"CheckMailto": false})
	tExpectIssueCount(t, hT, 0)
}

func TestTelValid(t *testing.T) {
	// ignores valid tel links
	hT := tTestFile("fixtures/links/tel_link.html")
	tExpectIssueCount(t, hT, 0)
}

func TestTelBlank(t *testing.T) {
	// fails for blank tel links
	hT := tTestFile("fixtures/links/blank_tel_link.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "tel is empty", 1)
}

func TestTelBlankIgnore(t *testing.T) {
	// fails for broken internal links
	hT := tTestFileOpts("fixtures/links/blank_tel_link.html",
		map[string]interface{}{"CheckTel": false})
	tExpectIssueCount(t, hT, 0)
}

func TestLinkHrefValid(t *testing.T) {
	// works for valid href within link elements
	hT := tTestFile("fixtures/links/head_link_href.html")
	tExpectIssueCount(t, hT, 0)
}

func TestLinkHrefBlank(t *testing.T) {
	// fails for empty href within link elements
	hT := tTestFile("fixtures/links/head_link_href_empty.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "<link> href blank", 1)
}

func TestLinkHrefBlankIgnore(t *testing.T) {
	// works for empty href within link elements when ignoring
	hT := tTestFileOpts("fixtures/links/head_link_href_empty.html",
		map[string]interface{}{"IgnoreEmptyHref": true})
	tExpectIssueCount(t, hT, 0)
}

func TestLinkHrefAbsent(t *testing.T) {
	// fails for absent href within link elements
	hT := tTestFile("fixtures/links/head_link_href_absent.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "<link> missing href", 1)
}

func TestLinkHrefBrokenCanonicalDefault(t *testing.T) {
	// works for valid href within link elements
	hT := tTestFileOpts("fixtures/links/brokenCanonicalLink.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 0)
}

func TestLinkHrefBrokenCanonicalOption(t *testing.T) {
	// works for valid href within link elements
	hT := tTestFileOpts("fixtures/links/brokenCanonicalLink.html",
		map[string]interface{}{"IgnoreCanonicalBrokenLinks": false, "VCREnable": true})
	tExpectIssueCount(t, hT, 1)
}

func TestLinkRelDnsPrefetch(t *testing.T) {
	// ignores links with rel="dns-prefetch"
	hT := tTestFile("fixtures/links/link-rel-dns-prefetch.html")
	tExpectIssueCount(t, hT, 0)
}

func TestLinkRelPreconnect(t *testing.T) {
	// ignores links with rel="preconnect"
	hT := tTestFile("fixtures/links/link_rel_preconnect.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorPre(t *testing.T) {
	// catches broken links when inside pre or code tags
	hT := tTestFileOpts("fixtures/links/anchors_in_pre.html",
		map[string]interface{}{"VCREnable": true})
	tExpectIssueCount(t, hT, 2)
	tExpectIssue(t, hT, "no such host", 2)
}

func TestLinkPre(t *testing.T) {
	// catches broken links when inside pre or code tags
	hT := tTestFile("fixtures/links/links_in_pre.html")
	tExpectIssueCount(t, hT, 2)
	tExpectIssue(t, hT, "<link> href blank", 2)
}

func TestAnchorHashQueryBroken(t *testing.T) {
	// fails for broken hash with query
	t.Skip("Not yet dealt with")
	hT := tTestFile("fixtures/links/broken_hash_with_query.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "PLACEHOLDER", 99)
}

func TestAnchorHashSelf(t *testing.T) {
	// works for hash referring to itself
	hT := tTestFile("fixtures/links/hashReferringToSelf.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorNameIgnore(t *testing.T) {
	// ignores placeholder with name
	hT := tTestFile("fixtures/links/placeholder_with_name.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorIdIgnore(t *testing.T) {
	// ignores placeholder with id
	hT := tTestFile("fixtures/links/placeholder_with_id.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorIdEmpty(t *testing.T) {
	// fails for placeholder with empty id
	// TODO: Should we only fail here if missing href?
	t.Skip("Not yet implemented")
	hT := tTestFile("fixtures/links/placeholder_with_empty_id.html")
	tExpectIssueCount(t, hT, 1)
	tExpectIssue(t, hT, "anchor with empty id", 99)
}

func TestAnchorOtherProtocols(t *testing.T) {
	// ignores non-hypertext protocols
	hT := tTestFile("fixtures/links/other_protocols.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorBlankHTML5(t *testing.T) {
	// does not expect href for anchors in HTML5
	hT := tTestFile("fixtures/links/blank_href_html5.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorBlankHTML4(t *testing.T) {
	// does expect href for anchors in non-HTML5
	t.Skip("Not yet implemented")
	hT1 := tTestFile("fixtures/links/blank_href_html4.html")
	tExpectIssueCount(t, hT1, 1)
	hT2 := tTestFile("fixtures/links/blank_href_htmlunknown.html")
	tExpectIssueCount(t, hT2, 1)
}

func TestAnchorInternalAbsolute(t *testing.T) {
	// works for internal absolute links
	hT := tTestFile("fixtures/links/absoluteLinks.html")
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalAbsoluteRootPublication(t *testing.T) {
	// works for internal absolute links when site is published to root of domain
	hT := tTestFileOpts("fixtures/links/absoluteLinksRootPublication.html",
		map[string]interface{}{"BaseURL": "http://example.com"})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalAbsoluteFolderPublication(t *testing.T) {
	// works for internal absolute links when site is published to a folder
	hT := tTestFileOpts("fixtures/links/absoluteLinksFolderPublication.html",
		map[string]interface{}{"BaseURL": "http://www.example.com/blog"})
	tExpectIssueCount(t, hT, 0)
}

func TestAnchorInternalBrokenAbsoluteFolderPublication(t *testing.T) {
	// works for missing internal absolute links when site is published to a folder
	hT := tTestFileOpts("fixtures/links/absoluteBrokenLinksFolderPublication.html",
		map[string]interface{}{"BaseURL": "http://www.example.com/blog"})
	tExpectIssueCount(t, hT, 2)
}

func TestAnchorExternalLinksToInternalFiles(t *testing.T) {
	// works for external links that reference internal files
	hT := tTestFileOpts("fixtures/links/externalLinksToInternalFiles.html",
		map[string]interface{}{
			"BaseURL":                       "http://www.example.com/blog",
			"CheckSelfReferencesAsInternal": true,
		})
	tExpectIssueCount(t, hT, 0)
}

func TestSelfSignedLink(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/selfSignedLink.html",
		map[string]interface{}{"IgnoreSSLVerify": false})
	tExpectIssueCount(t, hT, 1)
}

func TestSelfSignedLinkIgnoreSSLVerify(t *testing.T) {
	tSkipShortExternal(t)
	hT := tTestFileOpts("fixtures/links/selfSignedLink.html",
		map[string]interface{}{"IgnoreSSLVerify": true})
	tExpectIssueCount(t, hT, 0)
}

// Favicon

func TestFaviconDefaultMissing(t *testing.T) {
	// passes, by default, for missing favicon
	hT := tTestFile("fixtures/favicon/favicon_absent.html")
	tExpectIssue(t, hT, "favicon missing", 0)
}

func TestFaviconOptionMissing(t *testing.T) {
	// fails, when asked, for missing favicon
	hT := tTestFileOpts("fixtures/favicon/favicon_absent.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 1)
}

func TestFaviconOptionMissingApple(t *testing.T) {
	// fails, when asked, for present apple icon but missing favicon
	hT := tTestFileOpts("fixtures/favicon/favicon_absent_apple.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 1)
}

func TestFaviconOptionBroken(t *testing.T) {
	// fails for broken favicon
	hT := tTestFileOpts("fixtures/favicon/favicon_broken.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "target does not exist", 1)
}

func TestFaviconOptionBrokenIgnored(t *testing.T) {
	// fails with missing favicon for ignored icon link tag
	hT := tTestFileOpts("fixtures/favicon/favicon_broken_but_ignored.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 1)
}

func TestFaviconOptionPresent(t *testing.T) {
	// passes, when asked, for present favicon
	hT := tTestFileOpts("fixtures/favicon/favicon_present.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 0)
}

func TestFaviconOptionPresentShortcut(t *testing.T) {
	// passes, when asked, with present favicon with legacy rel="shortcut icon"
	hT := tTestFileOpts("fixtures/favicon/favicon_present_shortcut.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 0)
}

func TestFaviconOptionPresentButInBody(t *testing.T) {
	// fails when favicon isn't a first level child of <head>
	hT := tTestFileOpts("fixtures/favicon/favicon_present_but_in_body.html",
		map[string]interface{}{"CheckFavicon": true})
	tExpectIssue(t, hT, "favicon missing", 1)
}

// Benchmarks

func BenchmarkManyExternalLinks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tTestFileOpts("fixtures/benchmarks/manyExternalLinks.html",
			map[string]interface{}{"LogLevel": issues.LevelNone})
	}
}
