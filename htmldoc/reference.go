package htmldoc

import (
	"golang.org/x/net/html"
	"net/url"
	"path"
	"strings"
)

// Representation of the link between a document and a resource
type Reference struct {
	Document *Document  // Document node is in
	Node     *html.Node // Node reference was created from
	Path     string     // href/src taken verbatim from source
	URL      *url.URL   // URL object created from Path
}

// Create a new reference given a document, node and path. Generates the URL
// object.
func NewReference(document *Document, node *html.Node, path string) *Reference {
	// Clean path
	path = strings.TrimLeftFunc(path, invalidPrePostRune)
	path = strings.TrimRightFunc(path, invalidPrePostRune)
	// Create ref
	ref := Reference{
		Document: document,
		Node:     node,
		Path:     path,
	}
	// Parse and store parsed URL
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	ref.URL = u
	return &ref
}

// Returns the scheme of the reference. Uses URL.Scheme and adds "file" and
// "self" schemes for inter-file and intra-file references.
func (ref *Reference) Scheme() string {
	if strings.HasPrefix(ref.Path, "//") {
		// Could be http or https, we can handle https so prefer that
		// TODO Should we test both?
		return "https"
	}

	switch ref.URL.Scheme {
	case "http":
		return "http"
	case "https":
		return "https"
	case "":
		if len(ref.URL.Path) > 0 {
			return "file"
		} else {
			return "self"
		}
	case "mailto":
		return "mailto"
	case "tel":
		return "tel"
	}
	return "" // Unknown
}

// Proxy for URL.String but deals with other valid URL types, such as missing
// protocol URLs.
func (ref *Reference) URLString() string {
	// Format url for use in http.Get
	urlStr := ref.URL.String()
	if strings.HasPrefix(ref.Path, "//") {
		return "https:" + ref.URL.String()
	}
	return urlStr
}

// Is an internal absolute link
func (ref *Reference) IsInternalAbsolute() bool {
	return !strings.HasPrefix(ref.Path, "//") && strings.HasPrefix(ref.Path, "/")
}

// For internals, return a path to the referenced file relative to the
// 'site root'.
func (ref *Reference) RefSitePath() string {
	if ref.IsInternalAbsolute() {
		return ref.URL.Path
	} else {
		return path.Join(ref.Document.BasePath, ref.URL.Path)
	}
}

// Utilities

// Removes query string from given urlStr
func URLStripQueryString(urlStr string) string {
	return strings.Split(urlStr, "?")[0]
}
