package htmltest

import (
	"github.com/wjdp/htmltest/htmldoc"
	"log"
	"os"
	"path"
)

// Walk through the directory tree and pick .html files
func RecurseDirectory(dPath string) []htmldoc.Document {
	documents := make([]htmldoc.Document, 0)

	// Open dPath
	f, err := os.Open(makePath(dPath))
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
					Path:      fPath,
				})
			}
		}
	} else {
		log.Fatalf("%s isn't a directory", dPath)
	}

	return documents
}
