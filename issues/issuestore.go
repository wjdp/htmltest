package issues

import (
  "log"
  "fmt"
  "github.com/fatih/color"
  "github.com/wjdp/htmltest/doc"
)

const NONE int = 99
const ERROR int = 3
const WARNING int = 2
const INFO int = 1
const DEBUG int = 0

var LogLevel int;

type Issue struct {
  Level int
  Document *doc.Document
  Reference *doc.Reference
  Message string
}

var issueStore []Issue

func InitIssueStore() {
  issueStore = make([]Issue, 0)
}

func AddIssue(issue Issue) {
  issueStore = append(issueStore, issue)
  PrintIssue(issue)
}

func PrintIssue(issue Issue) {
  if issue.Level < LogLevel {
    return
  }

  var primary string
  if issue.Document != nil {
    primary = issue.Document.Path
  } else if issue.Reference != nil {
    primary = issue.Reference.Document.Path
  } else {
    primary = "<nil>"
  }

  var secondary string
  if issue.Reference != nil {
    secondary = issue.Reference.Path
  } else {
    secondary = "<nil>"
  }

  issueText := fmt.Sprintf("%v --- %v --> %v", issue.Message, primary, secondary)

  switch issue.Level {
  case ERROR:
    color.Red(issueText)
  case WARNING:
    color.Yellow(issueText)
  case INFO:
    color.Blue(issueText)
  case DEBUG:
    color.Magenta(issueText)
  }
}

func OutputIssues() {
  for _, issue := range issueStore {
    log.Print(issue)
  }
}

func IssueCount(level int) int {
  c := 0
  for _, issue := range issueStore {
    if issue.Level >= level { c += 1 }
  }
  return c
}

func IssueMatchCount(message string) int {
  // Return number of issues with matching message
  c := 0
  for _, issue := range issueStore {
    if issue.Message == message { c += 1 }
  }
  return c
}

func Issues() []Issue {
  return issueStore
}
