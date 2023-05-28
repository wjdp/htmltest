package htmltest

import (
	"testing"

	"github.com/daviddengcn/go-assert"
	"github.com/theunrepentantgeek/htmltest/output"
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

	hT, err := Test(userOpts)
	output.CheckErrorPanic(err)

	assert.Equals(t, "hT.opts.CheckExternal", hT.opts.CheckExternal, false)
	assert.Equals(t, "hT.opts.LogLevel", hT.opts.LogLevel, 1337)
	assert.Equals(t, "hT.opts.ExternalTimeout", hT.opts.ExternalTimeout,
		defaults["ExternalTimeout"])
	assert.Equals(t, "hT.opts.RedirectLimit", hT.opts.RedirectLimit,
		defaults["RedirectLimit"])
	assert.Equals(t, "defaults.RedirectLimit", defaults["RedirectLimit"], -1)
}

func TestInList(t *testing.T) {
	lst := []interface{}{"alpha", "bravo", "charlie"}
	assert.Equals(t, "alpha in lst", InList(lst, "alpha"), true)
	assert.Equals(t, "bravo in lst", InList(lst, "bravo"), true)
	assert.Equals(t, "charlie in lst", InList(lst, "charlie"), true)
	assert.Equals(t, "delta in lst", InList(lst, "delta"), false)
}

func TestIsURLIgnored(t *testing.T) {
	userOpts := map[string]interface{}{
		"IgnoreURLs": []interface{}{"google.com", "test.example.com",
			"library.com", "//\\w+.assetstore.info/lib/"},
		"NoRun": true,
	}

	hT, err := Test(userOpts)
	output.CheckErrorPanic(err)

	assert.IsTrue(t, "url ignored", hT.opts.isURLIgnored("https://google.com/?q=1234"))
	assert.IsTrue(t, "url ignored", hT.opts.isURLIgnored("https://test.example.com/"))
	assert.IsTrue(t, "url ignored", hT.opts.isURLIgnored("https://www.library.com/page"))
	assert.IsTrue(t, "url ignored", hT.opts.isURLIgnored("https://cdn.assetstore.info/lib/test.js"))
	assert.IsFalse(t, "url left alone", hT.opts.isURLIgnored("https://froogle.com/?q=1234"))
	assert.IsFalse(t, "url left alone", hT.opts.isURLIgnored("http://assetstore.info/lib/test.js"))
}

func TestIsInternalURLIgnored(t *testing.T) {
	userOpts := map[string]interface{}{
		"IgnoreInternalURLs": []interface{}{"/misc/js/script.js"},
		"NoRun":              true,
	}

	hT, err := Test(userOpts)
	output.CheckErrorPanic(err)

	assert.IsTrue(t, "url ignored", hT.opts.isInternalURLIgnored("/misc/js/script.js"))
	assert.IsFalse(t, "url left alone", hT.opts.isInternalURLIgnored("misc/js/script.js"))
	assert.IsFalse(t, "url left alone", hT.opts.isInternalURLIgnored("/misc/js/script"))
}

func TestMergeHTTPHeaders(t *testing.T) {
	userOpts := map[string]interface{}{
		"HTTPHeaders": map[interface{}]interface{}{
			"Range":  "bytes=0-10",
			"Accept": "*/*",
		},
		"NoRun": true,
	}

	hT, err := Test(userOpts)
	output.CheckErrorPanic(err)

	assert.Equals(t, "url ignored", hT.opts.HTTPHeaders["Range"], "bytes=0-10")
}
