package htmldoc

import (
	"golang.org/x/net/html"
	"os"
)

type Document struct {
	Path      string
	Directory string
	File      *os.File
	HTMLNode  *html.Node
}
