package issues

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestIssueStoreNew(t *testing.T) {
	iS := NewIssueStore(ERROR)
	assert.Equals(t, "IssueStore LogLevel", iS.LogLevel, ERROR)
}

func TestIssueStoreAdd(t *testing.T) {
	iS := NewIssueStore(NONE)
	issue := Issue{Level: ERROR, Message: "test"}
	iS.AddIssue(issue)
	assert.Equals(t, "issue count", iS.Count(ERROR), 1)
}

func TestIssueStoreCount(t *testing.T) {
	iS := NewIssueStore(NONE)
	iS.AddIssue(Issue{Level: ERROR, Message: "error"})
	iS.AddIssue(Issue{Level: WARNING, Message: "warn"})
	iS.AddIssue(Issue{Level: INFO, Message: "notice"})
	assert.Equals(t, "issue count", iS.Count(ERROR), 1)
	assert.Equals(t, "issue count", iS.Count(WARNING), 2)
	assert.Equals(t, "issue count", iS.Count(INFO), 3)
}

func TestIssueStoreMessageMatchCount(t *testing.T) {
	iS := NewIssueStore(NONE)
	iS.AddIssue(Issue{Level: ERROR, Message: "error one"})
	iS.AddIssue(Issue{Level: WARNING, Message: "error two"})
	iS.AddIssue(Issue{Level: INFO, Message: "notice"})
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("carrots"), 0)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("error"), 2)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("two"), 1)
	assert.Equals(t, "issue message match count",
		iS.MessageMatchCount("notice"), 1)
}
