// Package htmldoc : Provides local document interface for htmltest. Models a
// store of documents, individual documents and their internal and external
// references.
package htmldoc

import (
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/wjdp/htmltest/output"
)

// DocumentStore struct, store of Documents including Document discovery
type DocumentStore struct {
	BasePath           string               // Path, relative to cwd, the site is located in
	BaseURL            *url.URL             // Base URL of the site
	IgnorePatterns     []interface{}        // Regexes of directories to ignore
	Documents          []*Document          // All of the documents, used to iterate over
	DocumentPathMap    map[string]*Document // Maps slash separated paths to documents
	DocumentExtension  string               // File extension to look for
	DirectoryIndex     string               // What file is the index of the directory
	IgnoreTagAttribute string               // Attribute to ignore element and children if found on element
}

// NewDocumentStore : Create and return a new Document store.
func NewDocumentStore() DocumentStore {
	return DocumentStore{
		Documents:       make([]*Document, 0),
		DocumentPathMap: make(map[string]*Document),
	}
}

// AddDocument : Add a document to the document store.
func (dS *DocumentStore) AddDocument(doc *Document) {
	// Save reference to document to various data stores
	dS.Documents = append(dS.Documents, doc)
	dS.DocumentPathMap[doc.SitePath] = doc
	// Pass some vars on
	doc.ignoreTagAttribute = dS.IgnoreTagAttribute
}

// Discover : Discover all documents within DocumentStore.BasePath.
func (dS *DocumentStore) Discover() {
	dS.discoverRecurse(".")
}

// Does dir match one of the IgnorePatterns?
func (dS *DocumentStore) isDirIgnored(dir string) bool {
	for _, item := range dS.IgnorePatterns {
		if ok, _ := regexp.MatchString(item.(string), dir+"/"); ok {
			return true
		}
	}
	return false
}

// Recursive function to discover documents by walking the file tree
func (dS *DocumentStore) discoverRecurse(dPath string) {
	// Open directory to scan
	f, err := os.Open(path.Join(dS.BasePath, dPath))
	output.CheckErrorPanic(err)
	defer f.Close()

	// Get FileInfo of directory (scan it)
	fi, err := f.Stat()
	output.CheckErrorPanic(err)

	if fi.IsDir() { // Double check we're dealing with a directory
		// Read all FileInfo-s from directory, Readdir(count int)
		fis, err := f.Readdir(-1)
		output.CheckErrorPanic(err)

		// Iterate over contents of directory
		for _, fileinfo := range fis {
			fPath := path.Join(dPath, fileinfo.Name())
			if fileinfo.IsDir() {
				// If item is a dir, we delve deeper
				dS.discoverRecurse(fPath)
			} else if path.Ext(fileinfo.Name()) == dS.DocumentExtension {
				// If a file, create and save document
				newDoc := &Document{
					FilePath:   path.Join(dS.BasePath, fPath),
					SitePath:   fPath,
					BasePath:   dPath,
					IgnoreTest: dS.isDirIgnored(dPath),
				}
				newDoc.Init()
				dS.AddDocument(newDoc)
			}
		}
	} else { // It's a file, return single file
		panic("discoverRecurse encountered a file: " + dPath)
	}

}

// ResolvePath : Resolves internal absolute paths to documents.
func (dS *DocumentStore) ResolvePath(refPath string) (*Document, bool) {
	// Match root document
	if refPath == "/" {
		d0, b0 := dS.DocumentPathMap[dS.DirectoryIndex]
		return d0, b0
	}

	if refPath[0] == '/' && len(refPath) > 1 {
		// Is an absolute link, remove the leading slash for map lookup
		if dS.BaseURL == nil {
			// No base URL, so `/` means our root
			refPath = refPath[1:]
		} else {
			// We have a Base URL, so need to trip off the base path if present
			refPath = strings.TrimPrefix(refPath, dS.BaseURL.Path)

			// We want to end up with a relative path, so remove leading '/' if present
			// (This happens if BaseURL does *not* end in '/')
			refPath = strings.TrimPrefix(refPath, "/")
		}
	}

	// Try path as-is, path.ext
	d1, b1 := dS.DocumentPathMap[refPath]
	if b1 {
		// as-is worked, return that
		return d1, b1
	}

	// Try as a directory, path.ext/index.html
	d2, b2 := dS.DocumentPathMap[path.Join(refPath, dS.DirectoryIndex)]
	return d2, b2
}

// ResolveRef : Proxy to ResolvePath via ref.RefSitePath()
func (dS *DocumentStore) ResolveRef(ref *Reference) (*Document, bool) {
	return dS.ResolvePath(ref.RefSitePath())
}
