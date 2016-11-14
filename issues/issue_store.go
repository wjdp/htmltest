// htmltest issue store, provides a store and issue structs.
package issues

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

// Store of htmltest issues.
type IssueStore struct {
	LogLevel   int
	issues     []Issue
	writeMutex *sync.Mutex
	byteLog    []byte
}

// Create an issuestore, assigns defaults and returns.
func NewIssueStore(logLevel int) IssueStore {
	iS := IssueStore{LogLevel: logLevel}
	iS.issues = make([]Issue, 0)
	iS.writeMutex = &sync.Mutex{}
	iS.byteLog = make([]byte, 0)
	return iS
}

// Add an issue to the issue store, thread safe.
func (iS *IssueStore) AddIssue(issue Issue) {
	issue.store = iS // Set ref to issue store in issue
	iS.writeMutex.Lock()
	iS.issues = append(iS.issues, issue)
	issue.print(false)
	if issue.Level >= iS.LogLevel {
		// Build byte slice to write out at the end
		iS.byteLog = append(iS.byteLog, []byte(issue.text()+"\n")...)
	}
	iS.writeMutex.Unlock()
}

// Count the number of issues in the store at, or above, the given level.
func (iS *IssueStore) Count(level int) int {
	count := 0
	for _, issue := range iS.issues {
		if issue.Level >= level {
			count += 1
		}
	}
	return count
}

// Count the number of issues in the store containing the provided substr.
func (iS *IssueStore) MessageMatchCount(substr string) int {
	count := 0
	for _, issue := range iS.issues {
		if strings.Contains(issue.Message, substr) {
			count += 1
		}
	}
	return count
}

// Write the issue store to the given path, filtered by logLevel given in
// NewIssueStore.
func (iS *IssueStore) WriteLog(path string) {
	err := ioutil.WriteFile(path, iS.byteLog, 0644)
	if err != nil {
		panic(err)
	}
}

// Dump all issues to stdout, called by test helpers when issue asserts fail.
func (iS *IssueStore) DumpIssues(force bool) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<")
	for _, issue := range iS.issues {
		issue.print(force)
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
}
