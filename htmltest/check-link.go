package htmltest

import (
	"crypto/x509"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"github.com/wjdp/htmltest/output"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func (hT *HTMLTest) checkLink(document *htmldoc.Document, node *html.Node) {
	attrs := htmldoc.ExtractAttrs(node.Attr,
		[]string{"href", "rel"})

	// Check if favicon
	if htmldoc.AttrPresent(node.Attr, "rel") &&
		(attrs["rel"] == "icon" || attrs["rel"] == "shortcut icon") &&
		node.Parent.Data == "head" {
		document.State.FaviconPresent = true
	}

	// Ignore if rel=dns-prefetch, see #40. If we have more cases like this a hashable type should be created and
	// checked against.
	if attrs["rel"] == "dns-prefetch" {
		return
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
				Message:   "anchor without href",
				Reference: ref,
			})
			return
		case "link":
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "link tag missing href",
				Reference: ref,
			})
			return
		}
	}

	// Blank href
	if attrs["href"] == "" {
		if !hT.opts.IgnoreEmptyHref {
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelError,
				Message:   "href blank",
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
				Message:   "empty hash",
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

	// var refDoc *htmldoc.Document
	if len(ref.URL.Fragment) == 0 {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "missing hash",
			Reference: ref,
		})
	}

	if len(ref.URL.Path) > 0 {
		// internal
		refDoc, _ := hT.documentStore.ResolveRef(ref)
		if !refDoc.IsHashValid(ref.URL.Fragment) {
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

	// split off domain, check mx, fallback to A or AAAA if that fails
	var dnserr *net.DNSError
	var ok bool

	domain := strings.Split(ref.URL.Opaque, "@")[1]

	// loop over the current domain until we have a valid result or have exhausted all possibilities
	for domain != "" {
		// if a simple MX lookup works, we are done, continue
		if _, err := net.LookupMX(domain); err == nil {
			break // success, time to exit
		} else if dnserr, ok = err.(*net.DNSError); !ok || dnserr.Err != "no such host" {
			// this isn't an error we are expecting to see here
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelWarning,
				Message:   "unable to perform LookupMX, unknown error",
				Reference: ref,
			})
			return
		}

		// do we have to restart because of a CNAME
		if cname, err := net.LookupCNAME(domain); err == nil && cname != domain {
			// we have a valid CNAME, try with that. Loops return NXDOMAIN by default
			domain = cname
			continue

		} else if dnserr, ok = err.(*net.DNSError); !ok || dnserr.Err != "no such host" {
			// this isn't an error we are expecting to see here
			hT.issueStore.AddIssue(issues.Issue{
				Level:     issues.LevelWarning,
				Message:   "unable to perform LookupCNAME, unknown error",
				Reference: ref,
			})
			return
		}

		// an A or AAAA record here would be valid
		if _, err := net.LookupHost(domain); err == nil {
			break // its not ideal, but a valid A/AAAA record is acceptable for email
		} else {
			dnserr, ok = err.(*net.DNSError)
			if !ok || dnserr.Err != "no such host" {
				// we shouldn't see this here
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.LevelWarning,
					Message:   "unable to perform LookupHost, unknown error",
					Reference: ref,
				})
				return
			}

			if dnserr.Err == "no such host" {
				// represents NXDOMAIN or no records
				hT.issueStore.AddIssue(issues.Issue{
					Level:     issues.LevelError,
					Message:   "email domain could not be resolved correctly",
					Reference: ref,
				})
				return
			}
		}

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
