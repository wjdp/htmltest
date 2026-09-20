// Package issues : htmltest issue store, provides a store and issue structs.
package issues

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/output"
)

// IssueStore : store of htmltest issues.
type IssueStore struct {
	logLevel         int                 // Level of errors to report
	printImmediately bool                // Print issues when added
	issues           []*Issue            // Retained issues, see shouldRetain
	issuesByDoc      map[string][]*Issue // Retained issues by Document.SitePath
	droppedDebug     int                 // Issues dropped by shouldRetain, so Count stays exact
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

// shouldRetain : Should this issue be kept in the store after it has been
// counted and (possibly) printed?
//
// LevelDebug issues are emitted per link, so on a large site they run into the
// millions and dominate memory. They can only ever be shown when the store is
// actually reporting at debug level, so unless it is they are counted and
// discarded rather than retained for the whole run. Every other level is kept.
//
// Note logLevel is a *reporting* threshold, not a retention one: a store
// created with LevelNone still retains errors so they can be counted, listed
// and written out.
func (iS *IssueStore) shouldRetain(issue *Issue) bool {
	return issue.Level > LevelDebug || iS.logLevel <= LevelDebug
}

// AddIssue : Add an issue to the issue store, thread safe.
func (iS *IssueStore) AddIssue(issue Issue) {
	issue.store = iS // Set ref to issue store in issue

	// Capture the issue's descriptions and drop its Document/Reference
	// pointers before it is stored, so it cannot pin a parsed document tree.
	issue.resolve()

	iS.storeMutex.Lock()

	primary := issue.primary()

	if iS.shouldRetain(&issue) {
		iS.issues = append(iS.issues, &issue)
		iS.issuesByDoc[primary] = append(iS.issuesByDoc[primary], &issue)
	} else {
		// Still counted, so Count stays exact. A plain int is used rather than
		// a map so this behaves like the issues slice if the store is copied.
		iS.droppedDebug++
	}

	if iS.printImmediately || primary == textNil {
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
	iS.storeMutex.RLock()
	defer iS.storeMutex.RUnlock()

	count := 0
	for _, issue := range iS.issues {
		if issue.Level >= level {
			count++
		}
	}
	// Issues shouldRetain discarded are debug level, so they only register
	// against a request that reaches down that far.
	if level <= LevelDebug {
		count += iS.droppedDebug
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
// provided substr. Only considers retained issues, see shouldRetain.
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

// DumpIssues : Dump all retained issues to stdout, called by test helpers when
// issue asserts fail. See shouldRetain.
func (iS *IssueStore) DumpIssues(force bool) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<")
	for _, issue := range iS.issues {
		issue.print(force, "")
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
}
