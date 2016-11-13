package issues

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type IssueStore struct {
	LogLevel   int
	issues     []Issue
	writeMutex *sync.Mutex
	byteLog    []byte
}

func NewIssueStore(logLevel int) IssueStore {
	iS := IssueStore{LogLevel: logLevel}
	iS.issues = make([]Issue, 0)
	iS.writeMutex = &sync.Mutex{}
	iS.byteLog = make([]byte, 0)
	return iS
}

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

func (iS *IssueStore) WriteLog(path string) {
	err := ioutil.WriteFile(path, iS.byteLog, 0644)
	if err != nil {
		panic(err)
	}
}

func (iS *IssueStore) DumpIssues(force bool) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<")
	for _, issue := range iS.issues {
		issue.print(force)
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
}
