package htmltest

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
	"path"
	"reflect"
	"regexp"
	"strings"
)

// Options struct for htmltest, user and default options are merged and mapped
// into an instance of this struct.
type Options struct {
	DirectoryPath string
	FilePath      string

	CheckAnchors bool
	CheckLinks   bool
	CheckImages  bool
	CheckScripts bool
	CheckMeta    bool
	CheckGeneric bool

	CheckExternal     bool
	CheckInternal     bool
	CheckInternalHash bool
	CheckMailto       bool
	CheckTel          bool
	CheckFavicon      bool
	CheckMetaRefresh  bool

	EnforceHTTPS bool

	IgnoreURLs []interface{}
	IgnoreDirs []interface{}

	IgnoreInternalEmptyHash             bool
	IgnoreCanonicalBrokenLinks          bool
	IgnoreAltMissing                    bool
	IgnoreDirectoryMissingTrailingSlash bool
	IgnoreTagAttribute                  string

	TestFilesConcurrently    bool
	DocumentConcurrencyLimit int
	HTTPConcurrencyLimit     int

	LogLevel int
	LogSort  string

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

// DefaultOptions returns a map of default options.
func DefaultOptions() map[string]interface{} {
	// Specify defaults here
	return map[string]interface{}{
		"CheckAnchors": true,
		"CheckLinks":   true,
		"CheckImages":  true,
		"CheckScripts": true,
		"CheckMeta":    true,
		"CheckGeneric": true,

		"CheckExternal":     true,
		"CheckInternal":     true,
		"CheckInternalHash": true,
		"CheckMailto":       true,
		"CheckTel":          true,
		"CheckFavicon":      false,
		"CheckMetaRefresh":  true,

		"EnforceHTTPS": false,

		"IgnoreURLs": []interface{}{},
		"IgnoreDirs": []interface{}{},

		"IgnoreInternalEmptyHash":             false,
		"IgnoreCanonicalBrokenLinks":          true,
		"IgnoreAltMissing":                    false,
		"IgnoreDirectoryMissingTrailingSlash": false,
		"IgnoreTagAttribute":                  "data-proofer-ignore",

		"TestFilesConcurrently":    false,
		"DocumentConcurrencyLimit": 128,
		"HTTPConcurrencyLimit":     16,

		"LogLevel": issues.LevelWarning,
		"LogSort":  "document",

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

func (hT *HTMLTest) setOptions(optsUser map[string]interface{}) {
	// Merge user and default options, set Opts var
	optsMap := DefaultOptions()
	mergo.MergeWithOverwrite(&optsMap, optsUser)
	hT.opts = Options{}
	mergo.MapWithOverwrite(&hT.opts, optsMap)

	// If debug dump the options struct
	if hT.opts.LogLevel == issues.LevelDebug {
		s := reflect.ValueOf(&hT.opts).Elem()
		typeOfT := s.Type()

		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Printf("%d: %s %s = %v\n", i,
				typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
	}
}

// InList tests if key is in a slice/list.
func InList(list []string, key string) bool {
	for _, item := range list {
		if strings.Contains(key, item) {
			return true
		}
	}
	return false
}

// Is the given URL ignored by the current configuration
func (opts *Options) isURLIgnored(url string) bool {
	for _, item := range opts.IgnoreURLs {
		if ok, _ := regexp.MatchString(item.(string), url); ok {
			return true
		}
	}
	return false
}
