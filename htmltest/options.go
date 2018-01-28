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
	DirectoryPath  string
	DirectoryIndex string
	FilePath       string
	FileExtension  string

	CheckDoctype bool
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

	EnforceHTML5 bool
	EnforceHTTPS bool

	IgnoreURLs []string
	IgnoreDirs []string

	IgnoreInternalEmptyHash             bool
	IgnoreCanonicalBrokenLinks          bool
	IgnoreAltMissing                    bool
	IgnoreDirectoryMissingTrailingSlash bool
	IgnoreTagAttribute                  string

	HTTPHeaders map[interface{}]interface{}

	TestFilesConcurrently    bool
	DocumentConcurrencyLimit int
	HTTPConcurrencyLimit     int

	LogLevel int
	LogSort  string

	ExternalTimeout    int
	StripQueryString   bool
	StripQueryExcludes []string

	EnableCache     bool
	EnableLog       bool
	OutputDir       string
	OutputCacheFile string
	OutputLogFile   string
	CacheExpires    string // Accepts golang time period strings, hours (16h) is really only useful option

	// --- Internals below here ---
	NoRun     bool   // When true does not run tests, used to inspect state in unit tests
	VCREnable bool   // When true patches the govcr httpClient to mock network calls
	Version   string // Instigator should set this to a version string
}

// DefaultOptions returns a map of default options.
func DefaultOptions() map[string]interface{} {
	// Specify defaults here
	return map[string]interface{}{
		"DirectoryIndex": "index.html",
		"FileExtension":  ".html",

		"CheckDoctype": true,
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

		"EnforceHTML5": false,
		"EnforceHTTPS": false,

		"IgnoreURLs": []string{},
		"IgnoreDirs": []string{},

		"IgnoreInternalEmptyHash":             false,
		"IgnoreCanonicalBrokenLinks":          true,
		"IgnoreAltMissing":                    false,
		"IgnoreDirectoryMissingTrailingSlash": false,
		"IgnoreTagAttribute":                  "data-proofer-ignore",

		"HTTPHeaders": map[string]string{
			"Range":  "bytes=0-0", // If server supports prevents body being sent
			"Accept": "*/*",       // We accept all content types
		},

		"TestFilesConcurrently":    false,
		"DocumentConcurrencyLimit": 128,
		"HTTPConcurrencyLimit":     16,

		"LogLevel": issues.LevelWarning,
		"LogSort":  "document",

		"ExternalTimeout":    15,
		"StripQueryString":   true,
		"StripQueryExcludes": []string{"fonts.googleapis.com"},

		"EnableCache":     true,
		"EnableLog":       true,
		"OutputDir":       path.Join("tmp", ".htmltest"),
		"OutputCacheFile": "refcache.json",
		"OutputLogFile":   "htmltest.log",
		"CacheExpires":    "336h",

		"NoRun":     false,
		"VCREnable": false,
		"Version":   "dev",
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
		if ok, _ := regexp.MatchString(item, url); ok {
			return true
		}
	}
	return false
}
