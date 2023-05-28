package issues

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/theunrepentantgeek/htmltest/htmldoc"
)

const (
	// LevelNone : option to suppress output, actual error types follow
	LevelNone int = 99
	// LevelError : Fatal problems, presence of an error causes tests to fail
	LevelError int = 3
	// LevelWarning : An advisory, tests still pass
	LevelWarning int = 2
	// LevelInfo : Verbose information, normally hidden, not too noisy
	LevelInfo int = 1
	// LevelDebug : Debug output, normally hidden, very noisy
	LevelDebug int = 0
	// Text substitution when primary or secondary part of issue is nil
	textNil string = "<nil>"
)

// Issue struct representing a single issue with a document.
// Set all except Document and Reference, set one or the other.
type Issue struct {
	Level     int                // Level of the issue, use the consts at the top of this file
	Document  *htmldoc.Document  // Document this issue pertains to
	Reference *htmldoc.Reference // Reference this issue pertains to
	Message   string             // Error message, keep short
	store     *IssueStore        // Internal ref to the store this issue is owned by
}

// Textual description of the primary item in the issue
func (issue *Issue) primary() string {
	if issue.Document != nil {
		return issue.Document.SitePath
	} else if issue.Reference != nil && issue.Reference.Document != nil {
		return issue.Reference.Document.SitePath
	}
	return textNil
}

// Textual description of the secondary item in the issue
func (issue *Issue) secondary() string {
	if issue.Reference != nil {
		return issue.Reference.Path
	}
	return textNil
}

// Text to print
func (issue *Issue) text() string {
	pri := issue.primary()
	sec := issue.secondary()
	if pri != textNil || sec != textNil {
		return fmt.Sprintf("%v --- %v --> %v", issue.Message, issue.primary(),
			issue.secondary())
	}
	return issue.Message
}

// Print to stdout with optional colour (controlled by color.NoColor - see main())
func (issue *Issue) print(force bool, prefix string) {
	if (issue.Level < issue.store.logLevel) && !force {
		return
	}

	switch issue.Level {
	case LevelError:
		color.Set(color.FgRed)
	case LevelWarning:
		color.Set(color.FgYellow)
	case LevelInfo:
		color.Set(color.FgBlue)
	case LevelDebug:
		color.Set(color.FgMagenta)
	}

	fmt.Println(prefix + issue.text())

	color.Unset()
}
