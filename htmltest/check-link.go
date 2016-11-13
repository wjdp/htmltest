package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func (hT *HtmlTest) checkLink(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"href", "rel", hT.opts.IgnoreTagAttribute})

	// Ignore if data-proofer-ignore set
	if htmldoc.AttrPresent(node.Attr, hT.opts.IgnoreTagAttribute) {
		return
	}

	// Check if favicon
	if htmldoc.AttrPresent(node.Attr, "rel") &&
		(attrs["rel"] == "icon" || attrs["rel"] == "shortcut icon") &&
		node.Parent.Data == "head" {
		document.State.FaviconPresent = true
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["href"])

	// Check for missing href, fail for link nodes
	if !htmldoc.AttrPresent(node.Attr, "href") {
		switch node.Data {
		case "a":
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.DEBUG,
				Message:   "anchor without href",
				Reference: ref,
			})
			return
		case "link":
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "link tag missing href",
				Reference: ref,
			})
			return
		}
	}

	// Blank href
	if attrs["href"] == "" {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "href blank",
			Reference: ref,
		})
		return
	}

	// href="#"
	if attrs["href"] == "#" {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "empty hash",
			Reference: ref,
		})
		return
	}

	// Route reference check
	switch ref.Scheme() {
	case "http":
		if hT.opts.EnforceHTTPS {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "is not an HTTPS target",
				Reference: ref,
			})
		}
		hT.checkExternal(ref)
	case "https":
		hT.checkExternal(ref)
	case "file":
		hT.checkInternal(ref)
	case "self":
		hT.checkInternalHash(ref)
	case "mailto":
		hT.checkMailto(ref)
	case "tel":
		hT.checkTel(ref)
	}

	// TODO: Other schemes
	// What to do about unknown schemes, could be perfectly valid or a typo.
	// Perhaps show a warning, which can be suppressed per-scheme in options.
	// Preload with a couple of common ones, ftp &c.

}

func (hT *HtmlTest) checkExternal(ref *htmldoc.Reference) {
	if !hT.opts.CheckExternal {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "skipping external check",
			Reference: ref,
		})
		return
	}

	urlStr := ref.URLString()

	// Does this url match an url ignore rule?
	if hT.opts.IsURLIgnored(urlStr) {
		return
	}

	if hT.opts.StripQueryString && !InList(hT.opts.StripQueryExcludes, urlStr) {
		urlStr = htmldoc.URLStripQueryString(urlStr)
	}
	var statusCode int

	cR, isCached := hT.refCache.Get(urlStr)

	if isCached && statusCodeValid(cR.StatusCode) {
		// If we have a valid result in cache, use that
		statusCode = cR.StatusCode
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "from cache",
			Reference: ref,
		})
	} else {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "fresh",
			Reference: ref,
		})
		urlUrl, err := url.Parse(urlStr)
		req := &http.Request{
			Method: "GET",
			URL:    urlUrl,
			Header: map[string][]string{
				"Range": {"bytes=0-0"}, // If server supports prevents body being sent
			},
		}

		hT.httpChannel <- true // Add to http concurrency limiter

		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.INFO,
			Message:   "hitting",
			Reference: ref,
		})

		resp, err := hT.httpClient.Do(req)

		<-hT.httpChannel // Bump off http concurrency limiter

		if err != nil {
			if strings.Contains(err.Error(), "dial tcp") {
				// Remove long prefix
				prefix := "Get " + urlStr + ": dial tcp: lookup "
				cleanedMessage := strings.TrimPrefix(err.Error(), prefix)
				// Add error
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.ERROR,
					Message:   cleanedMessage,
					Reference: ref,
				})
				return
			}
			if strings.Contains(err.Error(), "Client.Timeout") {
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.ERROR,
					Message:   "request exceeded our ExternalTimeout",
					Reference: ref,
				})
				return
			}

			// Unhandled client error, return generic error
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   err.Error(),
				Reference: ref,
			})
			log.Println("Unhandled httpClient error:", err.Error())
			return
		}
		// Save cached result
		hT.refCache.Save(urlStr, resp.StatusCode)
		statusCode = resp.StatusCode
	}

	switch statusCode {
	case http.StatusOK:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	case http.StatusPartialContent:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	default:
		attrs := htmldoc.ExtractAttrs(ref.Node.Attr, []string{"rel"})
		if attrs["rel"] == "canonical" && hT.opts.IgnoreCanonicalBrokenLinks {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.WARNING,
				Message:   http.StatusText(statusCode) + " [rel=\"canonical\"]",
				Reference: ref,
			})
		} else {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   http.StatusText(statusCode),
				Reference: ref,
			})
		}
	}

	// TODO check a hash id exists in external page if present in reference (URL.Fragment)
}

func (hT *HtmlTest) checkInternal(ref *htmldoc.Reference) {
	if !hT.opts.CheckInternal {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "skipping internal check",
			Reference: ref,
		})
		return
	}

	// First lookup in document store,
	refDoc, refExists := hT.documentStore.ResolvePath(ref.AbsolutePath())

	if refExists {
		// If path doesn't end in slash and the resolved ref is an index.html, complain
		if ref.URL.Path[len(ref.URL.Path)-1] != '/' && path.Base(refDoc.SitePath) == hT.opts.DirectoryIndex {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "target is a directory, href lacks trailing slash",
				Reference: ref,
			})
		}
	} else {
		// If that fails attempt to lookup with filesystem, resolve a path and check
		refOsPath := path.Join(hT.opts.DirectoryPath, ref.AbsolutePath())
		hT.checkFile(ref, refOsPath)
	}

	if len(ref.URL.Fragment) > 0 {
		// Is also a hash link
		hT.checkInternalHash(ref)
	}
}

func (hT *HtmlTest) checkInternalHash(ref *htmldoc.Reference) {
	if !hT.opts.CheckInternalHash {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "skipping hash check",
			Reference: ref,
		})
		return
	}

	// var refDoc *htmldoc.Document
	if len(ref.URL.Fragment) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "missing hash",
			Reference: ref,
		})
	}

	if len(ref.URL.Path) > 0 {
		// internal
		refDoc, _ := hT.documentStore.ResolvePath(ref.AbsolutePath())
		if !refDoc.IsHashValid(ref.URL.Fragment) {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "hash does not exist",
				Reference: ref,
			})
		}
	} else {
		// self
		if !ref.Document.IsHashValid(ref.URL.Fragment) {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "hash does not exist",
				Reference: ref,
			})
		}
	}
}

func (hT *HtmlTest) checkFile(ref *htmldoc.Reference, absPath string) {
	f, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "target does not exist",
			Reference: ref,
		})
		return
	}
	checkErr(err) // Crash on other errors

	if f.IsDir() {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "target is a directory, no index",
			Reference: ref,
		})
	}
}

func (hT *HtmlTest) checkMailto(ref *htmldoc.Reference) {
	if !hT.opts.CheckMailto {
		return
	}
	if len(ref.URL.Opaque) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "mailto is empty",
			Reference: ref,
		})
		return
	}
	if !strings.Contains(ref.URL.Opaque, "@") {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "contains an invalid email address",
			Reference: ref,
		})
		return
	}
}

func (hT *HtmlTest) checkTel(ref *htmldoc.Reference) {
	if !hT.opts.CheckTel {
		return
	}
	if len(ref.URL.Opaque) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   "tel is empty",
			Reference: ref,
		})
		return
	}
}
