// Package issues : htmltest issue store, provides a store and issue structs.
package issues

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/theunrepentantgeek/htmltest/htmldoc"
	"github.com/theunrepentantgeek/htmltest/output"
	"io/ioutil"
	"strings"
	"sync"
)

// IssueStore : store of htmltest issues.
type IssueStore struct {
	logLevel         int                 // Level of errors to report
	printImmediately bool                // Print issues when added
	issues           []*Issue            // All issues
	issuesByDoc      map[string][]*Issue // Issues by Document.SitePath
	storeMutex       *sync.RWMutex       // Mutex to control access to stores
	byteLog          []byte              // Bytestream of issues, built when issues are added and written to disk at end
}

// NewIssueStore : Create an issuestore, assigns defaults and returns.
func NewIssueStore(logLevel int, printImmediately bool) IssueStore {
	iS := IssueStore{logLevel: logLevel, printImmediately: printImmediately}
	iS.issues = make([]*Issue, 0)
	iS.issuesByDoc = make(map[string][]*Issue)
	iS.storeMutex = &sync.RWMutex{}
	iS.byteLog = make([]byte, 0)
	return iS
}

// AddIssue : Add an issue to the issue store, thread safe.
func (iS *IssueStore) AddIssue(issue Issue) {
	issue.store = iS // Set ref to issue store in issue

	iS.storeMutex.Lock()

	iS.issues = append(iS.issues, &issue)
	iS.issuesByDoc[issue.primary()] = append(
		iS.issuesByDoc[issue.primary()], &issue)

	if iS.printImmediately || issue.primary() == textNil {
		issue.print(false, "")
	}
	if issue.Level >= iS.logLevel {
		// Build byte slice to write out at the end
		iS.byteLog = append(iS.byteLog, []byte(issue.text()+"\n")...)
	}

	iS.storeMutex.Unlock()
}

// Count : Counts the number of issues in the store at, or above, the given
// level.
func (iS *IssueStore) Count(level int) int {
	count := 0
	for _, issue := range iS.issues {
		if issue.Level >= level {
			count++
		}
	}
	return count
}

// CountByDoc : Count the number of issues in the store at, or above, the given
// level pertaining to the provided document. Thread safe.
func (iS *IssueStore) CountByDoc(level int, doc *htmldoc.Document) int {
	iS.storeMutex.RLock()
	count := 0
	for _, issue := range iS.issuesByDoc[doc.SitePath] {
		if issue.Level >= level {
			count++
		}
	}
	iS.storeMutex.RUnlock()
	return count
}

// MessageMatchCount : Count the number of issues in the store containing the
// provided substr.
func (iS *IssueStore) MessageMatchCount(substr string) int {
	count := 0
	for _, issue := range iS.issues {
		if strings.Contains(issue.Message, substr) {
			count++
		}
	}
	return count
}

// PrintDocumentIssues : Print issues pertaining to a single document, given
// that document's SitePath. Respects log level.
func (iS *IssueStore) PrintDocumentIssues(doc *htmldoc.Document) {
	if iS.CountByDoc(iS.logLevel, doc) == 0 {
		if iS.logLevel == LevelDebug {
			color.Set(color.FgMagenta)
			fmt.Println(doc.SitePath)
			color.Unset()
		}
		return
	}
	iS.storeMutex.RLock()
	fmt.Println(doc.SitePath)
	for _, issue := range iS.issuesByDoc[doc.SitePath] {
		issue.print(false, "  ")
	}
	iS.storeMutex.RUnlock()
}

// WriteLog : Write the issue store to the given path, filtered by logLevel
// given in NewIssueStore.
func (iS *IssueStore) WriteLog(path string) {
	err := ioutil.WriteFile(path, iS.byteLog, 0644)
	output.CheckErrorPanic(err)
}

// DumpIssues : Dump all issues to stdout, called by test helpers when issue
// asserts fail.
func (iS *IssueStore) DumpIssues(force bool) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<")
	for _, issue := range iS.issues {
		issue.print(force, "")
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
}
