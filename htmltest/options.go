package htmltest

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/wjdp/htmltest/issues"
)

// Options struct for htmltest, user and default options are merged and mapped
// into an instance of this struct.
type Options struct {
	DirectoryPath  string
	DirectoryIndex string
	FilePath       string
	FileExtension  string

	BaseURL string

	CheckDoctype bool
	CheckAnchors bool
	CheckLinks   bool
	CheckImages  bool
	CheckScripts bool
	CheckMeta    bool
	CheckGeneric bool

	CheckExternal                 bool
	CheckInternal                 bool
	CheckInternalHash             bool
	CheckMailto                   bool
	CheckTel                      bool
	CheckFavicon                  bool
	CheckMetaRefresh              bool
	CheckSelfReferencesAsInternal bool

	EnforceHTML5 bool
	EnforceHTTPS bool

	IgnoreURLs         []interface{}
	IgnoreInternalURLs []interface{}
	IgnoreHTTPS        []interface{}
	IgnoreDirs         []interface{}

	IgnoreInternalEmptyHash             bool
	IgnoreEmptyHref                     bool
	IgnoreCanonicalBrokenLinks          bool
	IgnoreExternalBrokenLinks           bool
	IgnoreAltMissing                    bool
	IgnoreAltEmpty                      bool
	IgnoreDirectoryMissingTrailingSlash bool
	IgnoreSSLVerify                     bool
	IgnoreTagAttribute                  string

	HTTPHeaders map[interface{}]interface{}

	TestFilesConcurrently    bool
	DocumentConcurrencyLimit int
	HTTPConcurrencyLimit     int

	LogLevel int
	LogSort  string

	ExternalTimeout    int
	RedirectLimit      int
	StripQueryString   bool
	StripQueryExcludes []interface{}

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

		"IgnoreURLs":         []interface{}{},
		"IgnoreInternalURLs": []interface{}{},
		"IgnoreHTTPS":        []interface{}{},
		"IgnoreDirs":         []interface{}{},

		"IgnoreInternalEmptyHash":             false,
		"IgnoreEmptyHref":                     false,
		"IgnoreCanonicalBrokenLinks":          true,
		"IgnoreExternalBrokenLinks":           false,
		"IgnoreAltMissing":                    false,
		"IgnoreAltEmpty":                      false,
		"IgnoreDirectoryMissingTrailingSlash": false,
		"IgnoreSSLVerify":                     false,
		"IgnoreTagAttribute":                  "data-proofer-ignore",

		"HTTPHeaders": map[interface{}]interface{}{
			"Range":  "bytes=0-0", // If server supports prevents body being sent
			"Accept": "*/*",       // We accept all content types
		},

		"TestFilesConcurrently":    false,
		"DocumentConcurrencyLimit": 128,
		"HTTPConcurrencyLimit":     16,

		"LogLevel": issues.LevelWarning,
		"LogSort":  "document",

		"ExternalTimeout":    15,
		"RedirectLimit":      -1, // resort to built-in default
		"StripQueryString":   true,
		"StripQueryExcludes": []interface{}{"fonts.googleapis.com"},

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
	mergo.Merge(&optsMap, optsUser, mergo.WithOverride)
	hT.opts = Options{}
	mergo.Map(&hT.opts, optsMap, mergo.WithOverride)

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
func InList(list []interface{}, key string) bool {
	for _, item := range list {
		if strings.Contains(key, fmt.Sprintf("%s", item)) {
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

// Is the given URL an insecure URL ignored by IgnoreHTTPS
func (opts *Options) isInsecureURLIgnored(url string) bool {
	for _, item := range opts.IgnoreHTTPS {
		if ok, _ := regexp.MatchString(item.(string), url); ok {
			return true
		}
	}
	return false
}

// Solve #168
// Is the given local URL ignored by the current configuration
func (opts *Options) isInternalURLIgnored(url string) bool {
	for _, item := range opts.IgnoreInternalURLs {
		if item.(string) == url {
			return true
		}
	}
	return false
}
