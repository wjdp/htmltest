package test

import (
  "log"
  "os"
  "time"
  "path"
  "sync"
  "golang.org/x/net/html"
  "net/http"
  "github.com/wjdp/htmltest/issues"
  "github.com/wjdp/htmltest/doc"
)

var Opts Options
var basePath string

var httpClient *http.Client

func init() {
  transport := &http.Transport{
    TLSNextProto: nil, // Disable HTTP/2, "write on closed buffer" errors
  }
  httpClient = &http.Client{
    // Durations are in nanoseconds
    Transport: transport,
    Timeout: time.Duration(Opts.ExternalTimeout * 1000000000),
  }
}

func SetBasePath(bPath string) {
  // TODO integrate with Setup
  basePath = bPath
}

func makePath(p string) string {
  return path.Join(basePath, p)
}

func Test(opts Options) {
  OptionsSetDefaults(&opts)
  issues.InitIssueStore()
  Opts = opts

  if opts.FilePath != "" {
    doc := doc.Document{
      Directory: opts.DirectoryPath,
      Path: opts.FilePath,
    }
    TestFile(&doc)
  } else if opts.DirectoryPath != "" {
    TestDirectory(opts)
  } else {
    log.Fatal("Neither file or directory path provided")
  }
}


func TestDirectory(opts Options) {
  issues.LogLevel = Opts.LogLevel

  log.Printf("github.com/wjdp/htmltest started on %s", basePath)

  files := RecurseDirectory("")
  TestFiles(files)
  // issues.OutputIssues()

  log.Printf("%d files checked", len(files))
}

// Walk through the directory tree and pick .html files
func RecurseDirectory(dPath string) []doc.Document {
  documents := make([]doc.Document, 0)

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
        documents = append(documents, doc.Document{
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

func TestFiles(documents []doc.Document) {

  if Opts.TestFilesConcurrently {
    var wg sync.WaitGroup
    for _, document := range documents {
      wg.Add(1)
      go func(document doc.Document) {
        defer wg.Done()
        TestFile(&document)
      }(document)
    }
    wg.Wait()
  } else {
    for _, document := range documents {
      TestFile(&document)
    }
  }
}

func TestFile(document *doc.Document) {
  // log.Println("testFile", document.Path)
  f, err := os.Open( path.Join(basePath, document.Path) )
  checkErr(err)
  defer f.Close()

  document.File = f

  parseHtml(document)
}

func parseHtml(document *doc.Document) {
  doc, err := html.Parse(document.File)
  checkErr(err)
  document.HTMLNode = doc
  parseNode(document, document.HTMLNode)
}

func parseNode(document *doc.Document, n *html.Node) {
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
