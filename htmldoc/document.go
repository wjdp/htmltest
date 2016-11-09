package htmldoc

import (
	"golang.org/x/net/html"
	"log"
	"os"
	"path"
)

type Document struct {
	FilePath  string // Relative to the shell session
	SitePath  string // Relative to the site root
	Directory string
	HTMLNode  *html.Node
}

func (doc *Document) Parse() {
	// Open, parse, and close document
	f, err := os.Open(doc.FilePath)
	checkErr(err)
	defer f.Close()

	htmlNode, err := html.Parse(f)
	checkErr(err)

	doc.HTMLNode = htmlNode
}

func DocumentsFromDir(path string) []Document {
	// Nice proxy for recurseDir
	return recurseDir(path, "")
}

func recurseDir(basePath string, dPath string) []Document {
	// Recursive function that returns all Document struts in a given
	// os directory.
	// basePath: the directory to scan
	// dPath: the subdirectory within basePath we're scanning

	documents := make([]Document, 0)

	// Open directory to scan
	f, err := os.Open(path.Join(basePath, dPath))
	checkErr(err)
	defer f.Close()

	// Get FileInfo of directory (scan it)
	fi, err := f.Stat()
	checkErr(err)

	if fi.IsDir() { // Double check we're dealing with a directory
		// Read all FileInfo-s from directory, Readdir(count int)
		fis, err := f.Readdir(-1)
		checkErr(err)

		// Iterate over contents of directory
		for _, fileinfo := range fis {
			fPath := path.Join(dPath, fileinfo.Name())
			if fileinfo.IsDir() {
				// If item is a dir, we need to iterate further, save returned documents
				documents = append(documents, recurseDir(basePath, fPath)...)
			} else if path.Ext(fileinfo.Name()) == ".html" || path.Ext(fileinfo.Name()) == ".htm" {
				// If a file, save to filename list
				documents = append(documents, Document{
					FilePath:  path.Join(basePath, fPath),
					SitePath:  fPath,
					Directory: dPath,
				})
			}
		}
	} else { // It's a file, fall over
		log.Fatalf("%s isn't a directory", dPath)
	}

	return documents
}
