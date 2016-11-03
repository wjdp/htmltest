package doc

import (
	"golang.org/x/net/html"
	"log"
	"net/url"
	"path"
	"strings"
)

type Reference struct {
	Document *Document
	Node     *html.Node
	Path     string
	URL      *url.URL
	Scheme   string
}

func NewReference(document *Document, node *html.Node, path string) *Reference {
	ref := Reference{
		Document: document,
		Node:     node,
		Path:     path,
	}
	u, err := url.Parse(path)
	if err != nil {
		log.Fatal(err)
	}
	ref.URL = u
	ref.Scheme = scheme(&ref)
	return &ref
}

func scheme(ref *Reference) string {
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
	return ""
}

// // Is an external link
// func IsExternal(ref *Reference) {

// }

// Is an internal absolute link
func IsAbsolute(ref *Reference) bool {
	return !strings.HasPrefix(ref.Path, "//") && strings.HasPrefix(ref.Path, "/")
}

func URLString(ref *Reference) string {
	urlStr := ref.URL.String()
	if strings.HasPrefix(ref.Path, "//") {
		return "https:" + ref.URL.String()
	}

	return urlStr
}

func URLStripQueryString(urlStr string) string {
	return strings.Split(urlStr, "?")[0]
}

// If internal, return a path to the referenced file relative to the 'site root'
// Strip shit off the end?

func AbsolutePath(ref *Reference) string {
	if IsAbsolute(ref) {
		return ref.URL.Path
	} else {
		return path.Join(ref.Document.Directory, ref.URL.Path)
	}
}
