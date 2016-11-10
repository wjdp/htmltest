package htmldoc

import (
	"golang.org/x/net/html"
	"net/url"
	"path"
	"strings"
)

type Reference struct {
	Document *Document
	Node     *html.Node
	Path     string
	URL      *url.URL
}

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

func (ref *Reference) URLString() string {
	// Format url for use in http.Get
	urlStr := ref.URL.String()
	if strings.HasPrefix(ref.Path, "//") {
		return "https:" + ref.URL.String()
	}
	return urlStr
}

func (ref *Reference) IsInternalAbsolute() bool {
	// Is an internal absolute link
	return !strings.HasPrefix(ref.Path, "//") && strings.HasPrefix(ref.Path, "/")
}

func (ref *Reference) AbsolutePath() string {
	// If external return unchanged
	if ref.Scheme() != "file" {
		return ref.URL.Path
	}
	// If internal, return a path to the referenced file relative to the 'site root'
	// Strip shit off the end?
	if ref.IsInternalAbsolute() {
		return ref.URL.Path
	} else {
		return path.Join(ref.Document.Directory, ref.URL.Path)
	}
}

// Utilities

func URLStripQueryString(urlStr string) string {
	return strings.Split(urlStr, "?")[0]
}
