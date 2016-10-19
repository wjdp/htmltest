package htmltest

import (
  "log"
  "os"
  "io"
  "path"
  // "sync"
  "golang.org/x/net/html"
  "issues"
)

var Opts Options
var basePath string

func init() {
  Opts = NewOptions()
}

func SetBasePath(bPath string) {
  basePath = bPath
}

func Go() {
  log.Printf("htmltest started on %s", basePath)

  filenames := RecurseFile("")
  TestFiles(filenames)
  issues.OutputIssues()

  log.Printf("%d files checked", len(filenames))
}

// Walk through the directory tree and pick .html files
func RecurseFile(dPath string) []string {
  filenames := make([]string, 0)

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
        // If item is a dir, we need to iterate further, save returned filenames
        filenames = append(filenames, RecurseFile(fPath)...)
      } else if path.Ext(fileinfo.Name()) == ".html" {
        // If a file, save to filename list
        filenames = append(filenames, fPath)
      }
    }
  } else {
    log.Fatalf("%s isn't a directory", dPath)
  }

  return filenames
}

func checkErr(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func TestFiles(filenames []string) {
  // var wg sync.WaitGroup
  for _, filename := range filenames {
    // wg.Add(1)
    // go func(filename string) {
    //   defer wg.Done()
    //   testFile(filename)
    // }(filename)
    testFile(filename)
  }
}

func testFile(fPath string) {
  f, err := os.Open( path.Join(basePath, fPath) )
  checkErr(err)
  defer f.Close()

  parseHtml(fPath, f)
}

func parseHtml(fPath string, r io.Reader) {
  doc, err := html.Parse(r)
  checkErr(err)
  parseNode(fPath, doc)
}

func parseNode(fPath string, n *html.Node) {
  if n.Type == html.ElementNode {
    switch n.Data {
    case "a":
      CheckLink(fPath, n)
    case "img":
      CheckImg(n)
    case "link":
      CheckLink(fPath, n)
    case "script":
      CheckScript(n)
    }
  }
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    parseNode(fPath, c)
  }
}
