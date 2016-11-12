package htmldoc

import (
	"os"
	"path"
	"regexp"
)

type DocumentStore struct {
	BasePath          string               // Path, relative to cwd, the site is located in
	IgnorePatterns    []interface{}        // Regexes of directories to ignore
	Documents         []*Document          // All of the documents, used to iterate over
	DocumentPathMap   map[string]*Document // Maps slash separated paths to documents
	DocumentExtension string               // File extension to look for
}

func NewDocumentStore() DocumentStore {
	return DocumentStore{
		Documents:       make([]*Document, 0),
		DocumentPathMap: make(map[string]*Document),
	}
}

func (dS *DocumentStore) AddDocument(doc *Document) {
	// Save reference to document to various data stores
	dS.Documents = append(dS.Documents, doc)
	dS.DocumentPathMap[doc.SitePath] = doc
}

func (dS *DocumentStore) Discover() {
	// Find all documents in BasePath
	dS.discoverRecurse(".")
}

func (dS *DocumentStore) isDirIgnored(dir string) bool {
	// Does path dir match IgnorePatterns?
	for _, item := range dS.IgnorePatterns {
		if ok, _ := regexp.MatchString(item.(string), dir+"/"); ok {
			return true
		}
	}
	return false
}

func (dS *DocumentStore) discoverRecurse(dPath string) {
	// Recurse over relative path dPath, saves found documents to dS
	if dS.isDirIgnored(dPath) {
		return
	}

	// Open directory to scan
	f, err := os.Open(path.Join(dS.BasePath, dPath))
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
				// If item is a dir, we delve deeper
				dS.discoverRecurse(fPath)
			} else if path.Ext(fileinfo.Name()) == ".html" || path.Ext(fileinfo.Name()) == ".htm" {
				// If a file, create and save document
				dS.AddDocument(&Document{
					FilePath:  path.Join(dS.BasePath, fPath),
					SitePath:  fPath,
					Directory: dPath,
				})
			}
		}
	} else { // It's a file, return single file
		panic("discoverRecurse encountered a file: " + dPath)
	}

}

func (dS *DocumentStore) DocumentExists(sitePath string) bool {
	_, ok := dS.DocumentPathMap[sitePath]
	return ok
}
