package test

import (
  "testing"
  "fmt"
  "path"
  "github.com/wjdp/htmltest/issues"
)

var t_LogLevel int = issues.WARNING

func t_expectIssue(t *testing.T, message string, expected int) {
  c := issues.IssueMatchCount(message)
  if c != expected {
    t.Error("expected issue", message, "count", expected, "!=", c)
    issues.OutputIssues()
  }
}

func t_expectIssueCount(t *testing.T, expected int) {
  c := issues.IssueCount(issues.WARNING)
  if c != expected {
    t.Error("expected", expected, "issues,", c, "found")
    issues.OutputIssues()
  }
}

func t_testFile(filename string) {
  opts := map[string]interface{}{
    "DirectoryPath": path.Dir(filename),
    "FilePath": path.Base(filename),
    "LogLevel": t_LogLevel,
  }
  Test(opts)
}

func t_testDirectory(filename string) {
  opts := map[string]interface{}{
    "DirectoryPath": path.Dir(filename),
    "LogLevel": t_LogLevel,
  }
  Test(opts)
}


func ExampleHelloWorld() {
  fmt.Println("Hello World")
  // Output:
  // Hello World
}

func TestAnchorMissingHref(t *testing.T) {
  // fails for link with no href
  t_testFile("fixtures/links/missingLinkHref.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "href blank", 1)
}

func TestAnchorIgnorable(t *testing.T) {
  // ignores links marked as ignore data-proofer-ignore
  t_testFile("fixtures/links/ignorableLinks.html")
  t_expectIssueCount(t, 0)
}



func TestExternalLinkBroken(t *testing.T) {
  // fails for broken external links
  t_testFile("fixtures/links/brokenLinkExternal.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "no such host", 1)
}

func TestExternalLinkIgnore(t *testing.T) {
  // ignores external links when asked
  filename := "fixtures/links/brokenLinkExternal.html"
  opts := map[string]interface{}{
    "DirectoryPath": path.Dir(filename),
    "FilePath": path.Base(filename),
    "LogLevel": t_LogLevel,
    "CheckExternal": false,
  }
  Test(opts)
  t_expectIssueCount(t, 0)
}

func TestExternalHashBrokenDefault(t *testing.T) {
  // passes for broken external hashes by default
  t_testFile("fixtures/links/brokenHashOnTheWeb.html")
  t_expectIssueCount(t, 0)
}

func TestExternalHashBrokenOption(t *testing.T) {
  // fails for broken external hashes when asked
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/brokenHashOnTheWeb.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "no such hash", 1)
}

func TestExternalCache(t *testing.T) {
  // does not check links with parameters multiple times
  // TODO check cache is being checked
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/check_just_once.html")
  t_expectIssueCount(t, 0)
}

func TestExternalHrefMalformed(t *testing.T) {
  // does not explode on bad external links in files
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/bad_external_links.html")
  t_expectIssueCount(t, 0) // TODO
}

func TestExternalInsecureDefault(t *testing.T) {
  // passes for non-HTTPS links when not asked
  t_testFile("fixtures/links/non_https.html")
  t_expectIssueCount(t, 0)
}

func TestExternalInsecureOption(t *testing.T) {
  // fails for non-HTTPS links when asked
  filename := "fixtures/links/non_https.html"
  opts := map[string]interface{}{
    "DirectoryPath": path.Dir(filename),
    "FilePath": path.Base(filename),
    "LogLevel": t_LogLevel,
    "EnforceHTTPS": true,
  }
  Test(opts)
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "is not an HTTPS link", 1)
}

func TestExternalHrefIP(t *testing.T) {
  // fails for broken IP address links
  t_testFile("fixtures/links/ip_href.html")
  t_expectIssueCount(t, 2)
  t_expectIssue(t, "request timed out", 2)
}

func TestExternalFollowRedirects(t *testing.T) {
  // should follow redirects
  t.Skip("Need new link, times out")
  t_testFile("fixtures/links/linkWithRedirect.html")
  t_expectIssueCount(t, 0)
}

func TestExternalFollowRedirectsDisabled(t *testing.T) {
  // fails on redirects if not following
  t.Skip("Not yet implemented, need new link, times out")
  t_testFile("fixtures/links/linkWithRedirect.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestExternalHTTPS(t *testing.T) {
  // should understand https
  t_testFile("fixtures/links/linkWithHttps.html")
  t_expectIssueCount(t, 0)
}

func TestExternalMissingProtocolValid(t *testing.T) {
  // works for valid links missing the protocol
  t_testFile("fixtures/links/link_missing_protocol_valid.html")
  t_expectIssueCount(t, 0)
}

func TestExternalMissingProtocolInvalid(t *testing.T) {
  // fails for invalid links missing the protocol
  t_testFile("fixtures/links/link_missing_protocol_invalid.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "no such host", 1)
}

func TestExternalHrefPipes(t *testing.T) {
  // works for pipes in the URL
  t_testFile("fixtures/links/escape_pipes.html")
  t_expectIssueCount(t, 0)
}



func TestInternalBroken(t *testing.T) {
  // fails for broken internal links
  t_testFile("fixtures/links/brokenLinkInternal.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "target does not exist", 1)
}

func TestInternalRelativeLinksBase(t *testing.T) {
  // passes for relative links with a base
  t.Skip("Broken, ones does not exist, third back operation to base not supported")
  t_testFile("fixtures/links/relativeLinksWithBase.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHashBroken(t *testing.T) {
  // fails for broken internal hash
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/brokenHashInternal.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestDirectoryRootResolve(t *testing.T) {
  // properly resolves implicit /index.html in link paths
  t_testFile("fixtures/links/linkToFolder.html")
  t_expectIssueCount(t, 0)
}

func TestDirectoryCustomRoot(t *testing.T) {
  // works for custom directory index file
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/link_pointing_to_directory.html")
  t_expectIssueCount(t, 0)
}

func TestDirectoryCustomRootBroken(t *testing.T) {
  // fails if custom directory index file doesn't exist
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/link_pointing_to_directory.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestDirectoryNoTrailingSlash(t *testing.T) {
  // fails for internal linking to a directory without trailing slash
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/link_directory_without_slash.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestDirectoryHtmlExtension(t *testing.T) {
  // works for custom directory index file
  t.Skip("Not yet implemented")
  t_testDirectory("fixtures/links/_site")
  t_expectIssueCount(t, 0)
}

func TestInternalRootLink(t *testing.T) {
  // properly checks links to root
  t_testFile("fixtures/links/rootLink/rootLink.html")
  t_expectIssueCount(t, 0)
}

func TestInternalRelativeLinks(t *testing.T) {
  // properly checks relative links
  t_testFile("fixtures/links/relativeLinks.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHrefNonstandardChars(t *testing.T) {
  // passes non-standard characters
  t_testFile("fixtures/links/non_standard_characters.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHrefUTF8(t *testing.T) {
  // passes for external UTF-8 links
  t_testFile("fixtures/links/utf8Link.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHrefUrlEncoded(t *testing.T) {
  // passes for urlencoded href
  t_testFile("fixtures/links/urlencoded-href.html")
  t_expectIssueCount(t, 0)
}

func TestErrorDuplication(t *testing.T) {
  // does not dupe errors
  t_testFile("fixtures/links/nodupe.html")
  t_expectIssueCount(t, 1)
}

func TestInternalDashedAttrs(t *testing.T) {
  // does not complain for files with attributes containing dashes
  t_testFile("fixtures/links/attributeWithDash.html")
  t_expectIssueCount(t, 0)
}

func TestInternalCaseMismatch(t *testing.T) {
  // does not complain for internal links with mismatched cases
  t_testFile("fixtures/links/ignores_cases.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHashDefault(t *testing.T) {
  // fails for # href when not asked
  t_testFile("fixtures/links/hash_href.html")
  t_expectIssue(t, "empty hash", 1)
  t_expectIssueCount(t, 1)
}

func TestInternalHashOption(t *testing.T) {
  // passes for # href when asked
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/hash_href.html")
  t_expectIssueCount(t, 0)
}

func TestInternalHashWeird(t *testing.T) {
  // works for internal links to weird encoding IDs
  t_testFile("fixtures/links/encodingLink.html")
  t_expectIssueCount(t, 0)
}



func TestMultipleProblems(t *testing.T) {
  // finds a mix of broken and unbroken links
  t.Skip("Only single problem, and an hash which is not yet supported.")
  // TODO make our own multiple problem file
  t_testFile("fixtures/links/multipleProblems.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestMailtoValid(t *testing.T) {
  // ignores valid mailto links
  t_testFile("fixtures/links/mailto_link.html")
  t_expectIssueCount(t, 0)
}

func TestMailtoBlank(t *testing.T) {
  // fails for blank mailto links
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/blank_mailto_link.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestMailtoInvalid(t *testing.T) {
  // fails for invalid mailto links
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/invalid_mailto_link.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestTelValid(t *testing.T) {
  // ignores valid tel links
  t_testFile("fixtures/links/tel_link.html")
  t_expectIssueCount(t, 0)
}

func TestTelBlank(t *testing.T) {
  // fails for blank tel links
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/blank_tel_link.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestJavascriptLinkIgnore(t *testing.T) {
  // ignores javascript links
  t_testFile("fixtures/links/javascript_link.html")
  t_expectIssueCount(t, 0)
}

func TestLinkHrefValid(t *testing.T) {
  // works for valid href within link elements
  t_testFile("fixtures/links/head_link_href.html")
  t_expectIssueCount(t, 0)
}

func TestLinkHrefBlank(t *testing.T) {
  // fails for empty href within link elements
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/head_link_href_empty.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestLinkHrefAbsent(t *testing.T) {
  // fails for absent href within link elements
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/head_link_href_absent.html")
  t_expectIssueCount(t, 99)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

// TODO invalid link href?

func TestPreAnchor(t *testing.T) {
  // works for broken anchors within pre
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/anchors_in_pre.html")
  t_expectIssueCount(t, 0)
}

func TestPreLink(t *testing.T) {
  // works for broken link within pre
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/links_in_pre.html")
  t_expectIssueCount(t, 0)
}

func TestHashQueryBroken(t *testing.T) {
  // fails for broken hash with query
  t.Skip("Not yet dealt with")
  t_testFile("fixtures/links/broken_hash_with_query.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "PLACEHOLDER", 99)
}

func TestHashSelf(t *testing.T) {
  // works for hash referring to itself
  t_testFile("fixtures/links/hashReferringToSelf.html")
  t_expectIssueCount(t, 0)
}

func TestAnchorNameIgnore(t *testing.T) {
  // ignores placeholder with name
  t_testFile("fixtures/links/placeholder_with_name.html")
  t_expectIssueCount(t, 0)
}

func TestAnchorIdIgnore(t *testing.T) {
  // ignores placeholder with id
  t_testFile("fixtures/links/placeholder_with_id.html")
  t_expectIssueCount(t, 0)
}

func TestAnchorIdEmpty(t *testing.T) {
  // fails for placeholder with empty id
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/placeholder_with_empty_id.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "anchor with empty id", 99)
}

func TestOtherProtocols(t *testing.T) {
  // ignores non-hypertext protocols
  t_testFile("fixtures/links/other_protocols.html")
  t_expectIssueCount(t, 0)
}

func TestAnchorBlankHTML5(t *testing.T) {
  // does not expect href for anchors in HTML5
  t_testFile("fixtures/links/blank_href_html5.html")
  t_expectIssueCount(t, 0)
}

func TestAnchorBlankHTML4(t *testing.T) {
  // does expect href for anchors in non-HTML5
  t.Skip("Not yet implemented")
  t_testFile("fixtures/links/blank_href_html4.html")
  t_expectIssueCount(t, 1)
  t_testFile("fixtures/links/blank_href_htmlunknown.html")
  t_expectIssueCount(t, 1)
}














func TestHTML5Page(t *testing.T) {
  // Page containing HTML5 tags
  t_testFile("fixtures/html/html5_tags.html")
  t_expectIssueCount(t, 0)
}

