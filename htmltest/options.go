package htmltest

import (
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"regexp"
	"strings"
)

type Options struct {
	DirectoryPath string
	FilePath      string

	CheckAnchors bool
	CheckLinks   bool
	CheckImages  bool
	CheckScripts bool

	CheckExternal     bool
	CheckInternal     bool
	CheckInternalHash bool
	CheckMailto       bool
	CheckTel          bool
	CheckFavicon      bool
	EnforceHTTPS      bool

	IgnoreURLs []interface{}
	IgnoreDirs []interface{}

	IgnoreCanonicalBrokenLinks          bool
	IgnoreAltMissing                    bool
	IgnoreDirectoryMissingTrailingSlash bool
	IgnoreTagAttribute                  string

	TestFilesConcurrently    bool
	DocumentConcurrencyLimit int
	HTTPConcurrencyLimit     int

	LogLevel int

	DirectoryIndex string

	ExternalTimeout    int
	StripQueryString   bool
	StripQueryExcludes []string

	EnableCache  bool
	EnableLog    bool
	ProgDir      string
	CacheFile    string
	LogFile      string
	CacheExpires string // Accepts golang time period strings, hours (16h) is really only useful option

	// --- Internals below here ---
	NoRun bool // When true does not run tests, used to inspect state in unit tests
}

func DefaultOptions() map[string]interface{} {
	// Specify defaults here
	return map[string]interface{}{
		"CheckAnchors": true,
		"CheckLinks":   true,
		"CheckImages":  true,
		"CheckScripts": true,

		"CheckExternal":     true,
		"CheckInternal":     true,
		"CheckInternalHash": true,
		"CheckMailto":       true,
		"CheckTel":          true,
		"CheckFavicon":      false,
		"EnforceHTTPS":      false,

		"IgnoreURLs": []interface{}{},
		"IgnoreDirs": []interface{}{},

		"IgnoreCanonicalBrokenLinks":          true,
		"IgnoreAltMissing":                    false,
		"IgnoreDirectoryMissingTrailingSlash": false,
		"IgnoreTagAttribute":                  "data-proofer-ignore",

		"TestFilesConcurrently":    false,
		"DocumentConcurrencyLimit": 128,
		"HTTPConcurrencyLimit":     16,

		"LogLevel": issues.WARNING,

		"DirectoryIndex": "index.html",

		"ExternalTimeout":    15,
		"StripQueryString":   true,
		"StripQueryExcludes": []string{"fonts.googleapis.com"},

		"EnableCache":  true,
		"EnableLog":    true,
		"ProgDir":      path.Join("tmp", ".htmltest"),
		"CacheFile":    "refcache.json",
		"LogFile":      "htmltest.log",
		"CacheExpires": "336h",

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

func (opts *Options) IsURLIgnored(url string) bool {
	for _, item := range opts.IgnoreURLs {
		if ok, _ := regexp.MatchString(item.(string), url); ok {
			return true
		}
	}
	return false
}
