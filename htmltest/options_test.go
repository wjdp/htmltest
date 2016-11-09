package htmltest

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	// Check DefaultOptions is returning something useful
	defaults := DefaultOptions()
	if _, ok := defaults["ExternalTimeout"]; !ok {
		t.Error("important bits missing from defaults")
	}
}

func TestSetOptions(t *testing.T) {
	// Check SetOptions assigns user options above default options
	defaults := DefaultOptions()
	userOpts := map[string]interface{}{
		"CheckExternal": false,
		"LogLevel":      1337,
		"NoRun":         true,
	}
	hT := Test(userOpts)
	assert.Equals(t, "hT.opts.CheckExternal", hT.opts.CheckExternal, false)
	assert.Equals(t, "hT.opts.LogLevel", hT.opts.LogLevel, 1337)
	assert.Equals(t, "hT.opts.ExternalTimeout", hT.opts.ExternalTimeout,
		defaults["ExternalTimeout"])
}

func TestInList(t *testing.T) {
	lst := []string{"alpha", "bravo", "charlie"}
	assert.Equals(t, "alpha in lst", InList(lst, "alpha"), true)
	assert.Equals(t, "bravo in lst", InList(lst, "bravo"), true)
	assert.Equals(t, "charlie in lst", InList(lst, "charlie"), true)
	assert.Equals(t, "delta in lst", InList(lst, "delta"), false)
}
