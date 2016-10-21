package htmltest

import (
  "log"
  "os"
  "path"
  "sync"
  "golang.org/x/net/html"
  "issues"
  "htmldoc"
)

var Opts Options
var basePath string

func init() {
  Opts = NewOptions()
}

func SetBasePath(bPath string) {
  basePath = bPath
}

func makePath(p string) string {
  return path.Join(basePath, p)
}

func Go() {
  issues.LogLevel = Opts.LogLevel

  log.Printf("htmltest started on %s", basePath)

  files := RecurseDirectory("")
  TestFiles(files)
  // issues.OutputIssues()

  log.Printf("%d files checked", len(files))
}

// Walk through the directory tree and pick .html files
func RecurseDirectory(dPath string) []htmldoc.Document {
  documents := make([]htmldoc.Document, 0)

  // Open dPath
  f, err := os.Open( path.Join(basePath, dPath) )
  checkErr(err)
  defer f.Close()

  // Get FileInfo of dPath
  fi, err := f.Stat()
  checkErr(err)

  if fi.IsDir() {
    // Read all FileInfo-s from dPath
    fis, err := f.Readdir(-1)
    checkErr(err)

    // Iterate over contents of dPath
    for _, fileinfo := range fis {
      fPath := path.Join(dPath, fileinfo.Name())
      if fileinfo.IsDir() {
        // If item is a dir, we need to iterate further, save returned documents
        documents = append(documents, RecurseDirectory(fPath)...)
      } else if path.Ext(fileinfo.Name()) == ".html" {
        // If a file, save to filename list
        documents = append(documents, htmldoc.Document{
          Directory: dPath,
          Path: fPath,
        })
      }
    }
  } else {
    log.Fatalf("%s isn't a directory", dPath)
  }

  return documents
}

func checkErr(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func TestFiles(documents []htmldoc.Document) {

  if Opts.TestFilesConcurrently {
    var wg sync.WaitGroup
    for _, document := range documents {
      wg.Add(1)
      go func(document htmldoc.Document) {
        defer wg.Done()
        testFile(&document)
      }(document)
    }
    wg.Wait()
  } else {
    for _, document := range documents {
      testFile(&document)
    }
  }
}

func testFile(document *htmldoc.Document) {
  // log.Println("testFile", document.Path)
  f, err := os.Open( path.Join(basePath, document.Path) )
  checkErr(err)
  defer f.Close()

  document.File = f

  parseHtml(document)
}

func parseHtml(document *htmldoc.Document) {
  doc, err := html.Parse(document.File)
  checkErr(err)
  document.HTMLNode = doc
  parseNode(document, document.HTMLNode)
}

func parseNode(document *htmldoc.Document, n *html.Node) {
  if n.Type == html.ElementNode {
    switch n.Data {
    case "a":
      CheckLink(document, n)
    case "img":
      CheckImg(n)
    case "link":
      CheckLink(document, n)
    case "script":
      CheckScript(n)
    }
  }
  // Iterate over children
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    parseNode(document, c)
  }
}
