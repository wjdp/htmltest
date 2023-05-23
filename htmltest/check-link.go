package htmltest

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/output"
	"golang.org/x/net/html"
)

// ignoredRels: List of rel values to ignore, dns-prefetch and preconnect are ignored as they are not links to be
//              followed rather telling browser we want something on that host, if the root of that host is not valid,
//              it's likely not a problem.
var ignoredRels = [...]string{"dns-prefetch", "preconnect"}

func (hT *HTMLTest) checkLink(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"href", "rel"})

	// Check if favicon
	if htmldoc.AttrPresent(node.Attr, "rel") &&
		(attrs["rel"] == "icon" || attrs["rel"] == "shortcut icon") &&
		node.Parent.Data == "head" {
		document.State.FaviconPresent = true
	}

	// If rel in IgnoredRels, ignore this link
	for _, rel := range ignoredRels {
		if attrs["rel"] == rel {
			return
		}
	}

	// Create reference
	ref, err := htmldoc.NewReference(document, node, attrs["href"])
	if err != nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:    issues.LevelError,
			Document: document,
			Message:  fmt.Sprintf("bad reference: %q", err),
		})
		return
	}

	// Check for missing href, fail for link nodes
	if !htmldoc.AttrPresent(node.Attr, "href") {
		switch node.Data {
		case "a":
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelDebug,
				Message:   "<a> without href",
				Reference: ref,
			})
			return
		case "link":
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "<link> missing href",
				Reference: ref,
			})
			return
		}
	}

	// Blank href
	if attrs["href"] == "" {
		if !hT.opts.IgnoreEmptyHref {
			var msg string = fmt.Sprintf("<%s> href blank", node.Data)
			if attrs["title"] != "" {
				msg = fmt.Sprintf("%s title=%q", msg, attrs["title"])
			}
			if ref.Node.FirstChild != nil {
				msg = fmt.Sprintf("%s body=%q", msg, ref.Node.FirstChild.Data)
			}
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   msg,
				Reference: ref,
			})
		}
		return
	}

	// href="#"
	if attrs["href"] == "#" {
		if hT.opts.CheckInternalHash && !hT.opts.IgnoreInternalEmptyHash {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   fmt.Sprintf("<%s> empty hash", node.Data),
				Reference: ref,
			})
		}
		return
	}

	// Route reference check
	switch ref.Scheme() {
	case "http":
		hT.enforceHTTPS(ref)
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

func (hT *HTMLTest) checkExternal(ref *htmldoc.Reference) {
	issueLevel := issues.LevelError
	if hT.opts.IgnoreExternalBrokenLinks {
		issueLevel = issues.LevelWarning
	}
	if !hT.opts.CheckExternal {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   "skipping external check",
			Reference: ref,
		})
		return
	}

	// Is this an external reference to a local file?
	if hT.opts.CheckSelfReferencesAsInternal && hT.documentStore.BaseURL != nil {

		if ref.URL.Host == hT.documentStore.BaseURL.Host && hT.documentStore.BaseURL.User == nil {
			// Convert to internal reference
			internalURL := *ref.URL
			internalURL.Scheme = ""
			internalURL.Host = ""

			internalRef := *ref
			internalRef.URL = &internalURL
			internalRef.Path = internalURL.String()

			hT.checkInternal(&internalRef)
			return
		}
	}

	urlStr := ref.URLString()

	// Does this url match an url ignore rule?
	if hT.opts.isURLIgnored(urlStr) {
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
			Level:     issues.LevelDebug,
			Message:   "from cache",
			Reference: ref,
		})
	} else {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   "fresh",
			Reference: ref,
		})

		// Build the request
		req, err := http.NewRequest("GET", urlStr, nil)
		// Only error NewRequest raises is if the url isn't valid, we have already checked it by this point so OK just
		// to panic if err != nil.
		output.CheckErrorPanic(err)

		// Set UA header
		req.Header.Set("User-Agent", "htmltest/"+hT.opts.Version)

		// Set headers from HTTPHeaders option
		for key, value := range hT.opts.HTTPHeaders {
			// Due to the way we're loading in config these keys and values are interface{}. In normal cases they are
			// strings, but could very easily be ints (side note: this isn't great, we'll fix this later, #73)
			req.Header.Set(fmt.Sprintf("%v", key), fmt.Sprintf("%v", value))
		}

		hT.httpChannel <- true // Add to http concurrency limiter

		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelInfo,
			Message:   "hitting",
			Reference: ref,
		})

		resp, err := hT.httpClient.Do(req)

		<-hT.httpChannel // Bump off http concurrency limiter

		if err != nil {
			if strings.Contains(err.Error(), "Client.Timeout") {
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issueLevel,
					Message:   "request exceeded our ExternalTimeout",
					Reference: ref,
				})
				return
			}

			if certErr, ok := err.(*url.Error).Err.(x509.UnknownAuthorityError); ok {
				err = validateCertChain(certErr.Cert)
				if err == nil {
					hT.issueStore.AddIssue(issues.Issue{
						Level:     issues.LevelWarning,
						Reference: ref,
						Message:   "incomplete certificate chain",
					})
					return
				}
			}

			// More generic, should be kept below more specific cases
			if strings.Contains(err.Error(), "dial tcp") {
				// Remove long prefix
				prefix := "Get " + urlStr + ": dial tcp: lookup "
				cleanedMessage := strings.TrimPrefix(err.Error(), prefix)
				// Add error
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issueLevel,
					Message:   cleanedMessage,
					Reference: ref,
				})
				return
			}

			// Unhandled client error, return generic error
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issueLevel,
				Message:   err.Error(),
				Reference: ref,
			})

			return
		}
		// Save cached result
		hT.refCache.Save(urlStr, resp.StatusCode)
		statusCode = resp.StatusCode
	}

	switch statusCode {
	case http.StatusOK:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	case http.StatusPartialContent:
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   http.StatusText(statusCode),
			Reference: ref,
		})
	default:
		attrs := htmldoc.ExtractAttrs(ref.Node.Attr, []string{"rel"})
		if attrs["rel"] == "canonical" && hT.opts.IgnoreCanonicalBrokenLinks {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelWarning,
				Message:   http.StatusText(statusCode) + " [rel=\"canonical\"]",
				Reference: ref,
			})
		} else {
			// Failed VCRed requests end up here with a status code of zero
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issueLevel,
				Message:   fmt.Sprintf("%s %d", "Non-OK status:", statusCode),
				Reference: ref,
			})
		}
	}

	// TODO check a hash id exists in external page if present in reference (URL.Fragment)
}

func (hT *HTMLTest) checkInternal(ref *htmldoc.Reference) {
	if !hT.opts.CheckInternal {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   "skipping internal check",
			Reference: ref,
		})
		return
	}

	urlStr := ref.URLString()

	// Does this internal url match either a standard URL ignore rule or internal
	// url ignore rule?
	if hT.opts.isInternalURLIgnored(urlStr) || hT.opts.isURLIgnored(urlStr) {
		return
	}

	// First lookup in document store,
	refDoc, refExists := hT.documentStore.ResolveRef(ref)

	if refExists {
		// If the resolved ref is an index.html and the path doesn't end in a
		// trailing slash (and isn't linking directly to the index), complain.
		if !hT.opts.IgnoreDirectoryMissingTrailingSlash && path.Base(refDoc.SitePath) == hT.opts.DirectoryIndex &&
			!strings.HasSuffix(ref.URL.Path, hT.opts.DirectoryIndex) && !strings.HasSuffix(ref.URL.Path, "/") {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "target is a directory, href lacks trailing slash",
				Reference: ref,
			})
			refExists = false
		}
	} else {
		// If that fails attempt to lookup with filesystem, resolve a path and check
		refOsPath := path.Join(hT.opts.DirectoryPath, ref.RefSitePath())
		refExists = hT.checkFile(ref, refOsPath)
	}

	if refExists && len(ref.URL.Fragment) > 0 {
		// Is also a hash link
		hT.checkInternalHash(ref)
	}
}

func (hT *HTMLTest) checkInternalHash(ref *htmldoc.Reference) {
	if !hT.opts.CheckInternalHash {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelDebug,
			Message:   "skipping hash check",
			Reference: ref,
		})
		return
	}

	if len(ref.URL.Fragment) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "missing hash",
			Reference: ref,
		})
	}

	if len(ref.URL.Path) > 0 {
		// internal
		refDoc, ok := hT.documentStore.ResolveRef(ref)

		if !ok || !refDoc.IsHashValid(ref.URL.Fragment) {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "hash does not exist",
				Reference: ref,
			})
		}
	} else {
		// self
		if !ref.Document.IsHashValid(ref.URL.Fragment) {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "hash does not exist",
				Reference: ref,
			})
		}
	}
}

func (hT *HTMLTest) checkFile(ref *htmldoc.Reference, absPath string) bool {
	f, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "target does not exist",
			Reference: ref,
		})
		return false
	}
	output.CheckErrorPanic(err)

	if f.IsDir() {
		f, err = os.Stat(path.Join(absPath, hT.opts.DirectoryIndex))
		if os.IsNotExist(err) {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "target is a directory, no index",
				Reference: ref,
			})
			return false
		}
	}
	return true
}

func (hT *HTMLTest) checkMailto(ref *htmldoc.Reference) {
	if !hT.opts.CheckMailto {
		return
	}
	if len(ref.URL.Opaque) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "mailto is empty",
			Reference: ref,
		})
		return
	}
	emailAddress, decodeErr := url.QueryUnescape(ref.URL.Opaque)
	if decodeErr != nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   fmt.Sprintf("cannot decode email (%s): '%s'", decodeErr, ref.URL.Opaque),
			Reference: ref,
		})
		return
	}
	formatErr := checkmail.ValidateFormat(emailAddress)
	if formatErr != nil {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   fmt.Sprintf("invalid email address (%s): '%s'", formatErr, emailAddress),
			Reference: ref,
		})
		return
	}
}

func (hT *HTMLTest) checkTel(ref *htmldoc.Reference) {
	if !hT.opts.CheckTel {
		return
	}
	if len(ref.URL.Opaque) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "tel is empty",
			Reference: ref,
		})
		return
	}
}
