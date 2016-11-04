package htmltest

import "testing"

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
	}
	SetOptions(userOpts)
	t_assertEqual(t, Opts.CheckExternal, false)
	t_assertEqual(t, Opts.LogLevel, 1337)
	t_assertEqual(t, Opts.ExternalTimeout, defaults["ExternalTimeout"])
}

func TestInList(t *testing.T) {
	lst := []string{"alpha", "bravo", "charlie"}
	t_assertEqual(t, InList(lst, "alpha"), true)
	t_assertEqual(t, InList(lst, "bravo"), true)
	t_assertEqual(t, InList(lst, "charlie"), true)
	t_assertEqual(t, InList(lst, "delta"), false)
}
