package issues

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/wjdp/htmltest/htmldoc"
)

const TEXT_NIL string = "<nil>"

const NONE int = 99
const ERROR int = 3
const WARNING int = 2
const INFO int = 1
const DEBUG int = 0

// Struct representing a single issue with a document.
// Set all except Document and Reference, set one or the other.
type Issue struct {
	Level     int                // Level of the issue, use the consts at the top of this file
	Document  *htmldoc.Document  // Document this issue pertains to
	Reference *htmldoc.Reference // Reference this issue pertains to
	Message   string             // Error message, keep short
	store     *IssueStore        // Internal ref to the store this issue is owned by
}

func (issue *Issue) primary() string {
	if issue.Document != nil {
		return issue.Document.SitePath
	} else if issue.Reference != nil && issue.Reference.Document != nil {
		return issue.Reference.Document.SitePath
	} else {
		return TEXT_NIL
	}
}

func (issue *Issue) secondary() string {
	if issue.Reference != nil {
		return issue.Reference.Path
	} else {
		return TEXT_NIL
	}
}

func (issue *Issue) text() string {
	return fmt.Sprintf("%v --- %v --> %v", issue.Message, issue.primary(),
		issue.secondary())
}

func (issue *Issue) print(force bool) {
	if (issue.Level < issue.store.LogLevel) && !force {
		return
	}

	switch issue.Level {
	case ERROR:
		color.Set(color.FgRed)
	case WARNING:
		color.Set(color.FgYellow)
	case INFO:
		color.Set(color.FgBlue)
	case DEBUG:
		color.Set(color.FgMagenta)
	}

	fmt.Println(issue.text())

	color.Unset()

}
