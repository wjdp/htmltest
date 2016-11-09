package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"strings"
)

type Options struct {
	DirectoryPath string
	FilePath      string

	CheckExternal bool
	CheckInternal bool
	CheckMailto   bool
	CheckTel      bool

	EnforceHTTPS bool

	IgnoreAlt bool

	TestFilesConcurrently    bool
	DocumentConcurrencyLimit int
	HTTPConcurrencyLimit     int

	LogLevel int

	DirectoryIndex string

	ExternalTimeout    int
	StripQueryString   bool
	StripQueryExcludes []string

	NoRun bool // When true does not run tests, used to inspect state in unit tests
}

func DefaultOptions() map[string]interface{} {
	// Specify defaults here
	return map[string]interface{}{
		"CheckExternal": true,
		"CheckInternal": true,
		"CheckMailto":   true,
		"CheckTel":      true,

		"EnforceHTTPS": false,

		"IgnoreAlt": false,

		"TestFilesConcurrently":    false,
		"DocumentConcurrencyLimit": 128,
		"HTTPConcurrencyLimit":     4,

		"LogLevel": issues.INFO,

		"DirectoryIndex": "index.html",

		"ExternalTimeout":    3,
		"StripQueryString":   true,
		"StripQueryExcludes": []string{"fonts.googleapis.com"},

		"NoRun": false,
	}
}

func (hT *HtmlTest) setOptions(optsUser map[string]interface{}) {
	// Merge user and default options, set Opts var
	optsMap := DefaultOptions()
	mergo.MergeWithOverwrite(&optsMap, optsUser)
	hT.opts = Options{}
	mergo.MapWithOverwrite(&hT.opts, optsMap)

}

func InList(list []string, key string) bool {
	for _, item := range list {
		if strings.Contains(key, item) {
			return true
		}
	}
	return false
}
