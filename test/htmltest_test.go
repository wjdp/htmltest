package test

import (
  "testing"
  "fmt"
  "path"
  "github.com/wjdp/htmltest/issues"
)

func t_expectIssue(t *testing.T, message string, expected int) {
  c := issues.IssueMatchCount(message)
  if c != expected {
    t.Error("expected issue", message, "count", expected, "!=", c)
    issues.OutputIssues()
  }
}

func t_expectIssueCount(t *testing.T, expected int) {
  c := issues.IssueCount(issues.WARNING)
  if c != expected {
    t.Error("expected", expected, "issues,", c, "found")
    issues.OutputIssues()
  }
}

func t_testFile(filename string) {
  opts := Options{
    DirectoryPath: path.Dir(filename),
    FilePath: path.Base(filename),
    // LogLevel: issues.NONE,
  }
  Test(opts)
}


func ExampleHelloWorld() {
  fmt.Println("Hello World")
  // Output:
  // Hello World
}

func TestBrokenExternalLinks(t *testing.T) {
  t_testFile("fixtures/links/brokenLinkExternal.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "no such host", 1)
}

func TestBrokenInternalLinks(t *testing.T) {
  t_testFile("fixtures/links/brokenLinkInternal.html")
  t_expectIssueCount(t, 1)
  t_expectIssue(t, "target does not exist", 1)
}

func TestHTML5Page(t *testing.T) {
  // Page containing HTML5 tags
  t_testFile("fixtures/html/html5_tags.html")
  t_expectIssueCount(t, 0)
}

