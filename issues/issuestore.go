package issues

import (
	"strings"
	"sync"
)

type IssueStore struct {
	LogLevel   int
	issues     []Issue
	writeMutex *sync.Mutex
}

func NewIssueStore(logLevel int) IssueStore {
	iS := IssueStore{LogLevel: logLevel}
	iS.issues = make([]Issue, 0)
	iS.writeMutex = &sync.Mutex{}
	return iS
}

func (iS *IssueStore) AddIssue(issue Issue) {
	issue.store = iS // Set ref to issue store in issue
	iS.writeMutex.Lock()
	iS.issues = append(iS.issues, issue)
	issue.print()
	iS.writeMutex.Unlock()
}

func (iS *IssueStore) Count(level int) int {
	count := 0
	for _, issue := range iS.issues {
		if issue.Level >= level {
			count += 1
		}
	}
	return count
}

func (iS *IssueStore) MessageMatchCount(message string) int {
	// Return number of issues with matching message
	count := 0
	for _, issue := range iS.issues {
		if strings.Contains(issue.Message, message) {
			count += 1
		}
	}
	return count
}
