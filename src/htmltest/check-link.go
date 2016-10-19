package htmltest

import (
  // "log"
  "os"
  "path"
  "strings"
  "golang.org/x/net/html"
  "net/url"
  "issues"
)

func CheckLink(fPath string, n *html.Node) {
  attrs := extractAttrs(n.Attr, []string{"href"})
  if _, ok := attrs["href"]; ok {
    nHref := attrs["href"]
    nUrl, err := url.Parse(nHref)
    checkErr(err)

    switch nUrl.Scheme {
    case "http":
      if Opts.EnforceHTTPS {
        issues.AddIssue(issues.Issue{
          Message: "is not an HTTPS link",
          Path: fPath,
          NUrl: nUrl,
        })
      }
      CheckExternal(fPath, nUrl)
    case "https":
      CheckExternal(fPath, nUrl)
    case "":
      CheckInternal(fPath, nHref, nUrl)
    case "mailto":
      CheckMailto(fPath, nUrl)
    case "tel":
      CheckMailto(fPath, nUrl)
    }

  } else {
    // Anchor without href, do... nothing?
  }
}

func CheckExternal(fPath string, nUrl *url.URL) {

}

func CheckInternal(fPath string, nHref string, nUrl *url.URL) {
  // TODO extract hashes
  // log.Print(nUrl)
  // log.Print(nUrl.Path)

  isAbsolute := strings.HasPrefix(nHref, "/")

  var filePath string

  if isAbsolute {
    filePath = path.Join(basePath, nUrl.Path)
  } else {
    filePath = path.Join(basePath, fPath)
  }

  if _, err := os.Stat( filePath ); os.IsNotExist(err) {
    issues.AddIssue(issues.Issue{
      Message: "does not exist",
      Path: fPath,
      NUrl: nUrl,
    })
  }
}

func CheckMailto(fPath string, nUrl *url.URL) {

}

func CheckTel(fPath string, nUrl *url.URL) {

}
