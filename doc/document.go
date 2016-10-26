package doc

import (
  "os"
  "golang.org/x/net/html"
)

type Document struct {
  Path string
  Directory string
  File *os.File
  HTMLNode *html.Node
}
