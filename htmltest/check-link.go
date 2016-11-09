package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/refcache"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func (hT *HtmlTest) checkLink(document *htmldoc.Document, node *html.Node) {
	attrs := extractAttrs(node.Attr, []string{"href", "rel", "data-proofer-ignore"})

	// Do not check canonical links
	if attrs["rel"] == "canonical" {
		return
	}
	// Ignore if data-proofer-ignore set
	if attrPresent(node.Attr, "data-proofer-ignore") {
		return
	}

	// Create reference
	ref := htmldoc.NewReference(document, node, attrs["href"])

	// Check for missing href, fail for link nodes
	if !attrPresent(node.Attr, "href") {
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
			Message:   "skipping",
			Reference: ref,
		})
		return
	}

	urlStr := ref.URLString()

	if hT.opts.StripQueryString && !InList(hT.opts.StripQueryExcludes, urlStr) {
		urlStr = htmldoc.URLStripQueryString(urlStr)
	}
	var statusCode int

	if refcache.CachedURLStatus(urlStr) != 0 {
		// If we have the result in cache, return that
		statusCode = refcache.CachedURLStatus(urlStr)
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.INFO,
			Message:   "from cache",
			Reference: ref,
		})
	} else {
		urlUrl, err := url.Parse(urlStr)
		req := &http.Request{
			Method: "GET",
			URL:    urlUrl,
			Header: map[string][]string{
				"Range": {"bytes=0-0"}, // If server supports prevents body being sent
			},
		}

		hT.httpChannel <- true // Add to http concurrency limiter
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
		refcache.SetCachedURLStatus(urlStr, resp.StatusCode)
		statusCode = resp.StatusCode
	}

	switch statusCode {
	case http.StatusOK:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.INFO,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	case http.StatusPartialContent:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.INFO,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	default:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.ERROR,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	}

	// TODO check a hash id exists in external page if present in reference (URL.Fragment)
}

func (hT *HtmlTest) checkInternal(ref *htmldoc.Reference) {
	if !hT.opts.CheckInternal {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.DEBUG,
			Message:   "skipping",
			Reference: ref,
		})
		return
	}
	// Resolve a filesystem path for reference
	refOsPath := path.Join(hT.opts.DirectoryPath, ref.AbsolutePath())
	hT.checkFile(ref, refOsPath)
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
		if !strings.HasSuffix(ref.Path, "/") {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.ERROR,
				Message:   "target is a directory, href lacks trailing slash",
				Reference: ref,
			})
			return
		}

		hT.checkFile(ref, path.Join(absPath, hT.opts.DirectoryIndex))
		return
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
