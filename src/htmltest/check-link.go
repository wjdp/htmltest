package htmltest

import (
  "log"
  "os"
  "path"
  // "strings"
  "golang.org/x/net/html"
  // "net/url"
  "net/http"
  "issues"
  "htmldoc"
)

func CheckLink(document *htmldoc.Document, node *html.Node) {
  attrs := extractAttrs(node.Attr, []string{"href"})
  if _, ok := attrs["href"]; ok {
    ref := htmldoc.NewReference(document, node, attrs["href"])
    switch ref.Scheme {
    case "http":
      if Opts.EnforceHTTPS {
        issues.AddIssue(issues.Issue{
          Level: issues.ERROR,
          Message: "is not an HTTPS link",
          Reference: ref,
          })
      }
      CheckExternal(ref)
    case "https":
      CheckExternal(ref)
    case "file":
      CheckInternal(ref)
    case "mailto":
    case "tel":

    }
  } else {
    issues.AddIssue(issues.Issue{
      Level: issues.DEBUG,
      Message: "anchor without href",
      Document: document,
    })
  }
}

func CheckExternal(ref *htmldoc.Reference) {
  if !Opts.CheckExternal {
    issues.AddIssue(issues.Issue{
      Level: issues.DEBUG,
      Message: "skipping",
      Reference: ref,
    })
    return
  }
  log.Println("Ext", htmldoc.URLString(ref))

  resp, err := http.Get(htmldoc.URLString(ref))

  if err != nil {
    issues.AddIssue(issues.Issue{
      Level: issues.ERROR,
      Message: err.Error(),
      Reference: ref,
    })
  }

  _ = resp

  // TODO check a hash id exists in external page if present in reference (URL.Fragment)
}

func CheckInternal(ref *htmldoc.Reference) {
  if !Opts.CheckInternal {
    issues.AddIssue(issues.Issue{
      Level: issues.DEBUG,
      Message: "skipping",
      Reference: ref,
    })
    return
  }
  // log.Println("CheckInternal", ref.Document.Path, htmldoc.AbsolutePath(ref))

  fPath := makePath(htmldoc.AbsolutePath(ref))
  CheckFile(ref, fPath)
}

func CheckFile(ref *htmldoc.Reference, fPath string) {
  f, err := os.Stat(fPath)
  if os.IsNotExist(err) {
    issues.AddIssue(issues.Issue{
      Level: issues.ERROR,
      Message: "target does not exist",
      Reference: ref,
    })
    return
  }
  checkErr(err) // Crash on other errors

  if f.IsDir() {
    issues.AddIssue(issues.Issue{
      Level: issues.DEBUG,
      Message: "target is a directory",
      Reference: ref,
    })
    CheckFile(ref, path.Join(fPath, Opts.DirectoryIndex))
    return
  }
}
