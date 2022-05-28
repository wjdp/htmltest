// Package issues : htmltest issue store, provides a store and issue structs.
package issues

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/output"
	"io/ioutil"
	"sort"
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

type IssueStats struct {
	// How many issues of each level
	TotalByLevel map[int]int
	// Collect errors against count
	ErrorsByMessage map[string]int
	// Collect warnings against count
	WarningsByMessage map[string]int
}

// GetIssueStats : Get stats on issues in the store.
func (iS *IssueStore) GetIssueStats() IssueStats {
	stats := IssueStats{TotalByLevel: make(map[int]int), ErrorsByMessage: make(map[string]int), WarningsByMessage: make(map[string]int)}
	for _, issue := range iS.issues {
		stats.TotalByLevel[issue.Level]++
		if issue.Level == LevelError {
			stats.ErrorsByMessage[issue.Message]++
		}
		if issue.Level == LevelWarning {
			stats.WarningsByMessage[issue.Message]++
		}
	}
	return stats
}

func formatMessageCounts(messageCounts map[string]int) string {
	keySlice := make([]string, 0)
	for key, _ := range messageCounts {
		keySlice = append(keySlice, key)
	}
	sort.Strings(keySlice)

	var s string
	for _, message := range keySlice {
		s += fmt.Sprintf("    %d \"%s\"\n", messageCounts[message], message)
	}
	return s
}

// FormatIssueStats : Return formatted stats on issues in the store.
func (iS *IssueStore) FormatIssueStats() string {
	formattedStats := ""
	stats := iS.GetIssueStats()
	if iS.logLevel <= LevelError {
		formattedStats += fmt.Sprintln("  Errors:  ", stats.TotalByLevel[LevelError])
	}
	if iS.logLevel <= LevelWarning {
		formattedStats += fmt.Sprintln("  Warnings:", stats.TotalByLevel[LevelWarning])
	}
	if iS.logLevel <= LevelInfo {
		formattedStats += fmt.Sprintln("  Infos:   ", stats.TotalByLevel[LevelInfo])
	}
	if iS.logLevel <= LevelDebug {
		formattedStats += fmt.Sprintln("  Debugs:  ", stats.TotalByLevel[LevelDebug])
	}
	if (iS.logLevel <= LevelError) && (len(stats.ErrorsByMessage) > 0) {
		formattedStats += fmt.Sprintln("  Errors by message:")
		formattedStats += formatMessageCounts(stats.ErrorsByMessage)
	}
	if (iS.logLevel <= LevelWarning) && (len(stats.WarningsByMessage) > 0) {
		formattedStats += fmt.Sprintln("  Warnings by message:")
		formattedStats += formatMessageCounts(stats.WarningsByMessage)
	}
	return formattedStats
}
