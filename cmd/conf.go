package cmd

import (
	"github.com/spf13/viper"
	"github.com/wjdp/htmltest/output"
	"fmt"
	"strings"
	"sort"
	"os"
)

type optsMap map[string]interface{}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func readConfig(configName string) {
	viper.SetConfigFile(configName)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if ucerr, ok := err.(viper.UnsupportedConfigError); ok {
			output.AbortWith(ucerr)
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			output.AbortWith("Config file", configName, "not found")
		} else {
			panic(err)
		}
		// TODO no such file or directory
	}
}

func viperToOpts() optsMap {
	// Reads from viper and outputs a dict that htmltest takes
	// We do not use viper in the htmltest package so this interface is required
	opts := optsMap{}

	// Add bools
	var boolOpts = []string{"CheckDoctype", "CheckAnchors", "CheckLinks", "CheckImages", "CheckScripts", "CheckMeta", "CheckGeneric", "CheckExternal", "CheckInternal", "CheckInternalHash", "CheckMailto", "CheckTel", "CheckFavicon", "CheckMetaRefresh", "EnforceHTML5", "EnforceHTTPS", "IgnoreInternalEmptyHash", "IgnoreCanonicalBrokenLinks", "IgnoreAltMissing", "IgnoreDirectoryMissingTrailingSlash", "TestFilesConcurrently", "StripQueryString", "EnableCache", "EnableLog"}

	for _, name := range boolOpts {
		if viper.IsSet(name) {
			opts[name] = viper.GetBool(name)
		}
	}

	// Add strings
	var stringOpts = []string{"DirectoryPath", "DirectoryIndex", "FilePath", "FileExtension", "IgnoreTagAttribute", "LogSort", "OutputDir", "OutputCacheFile", "OutputLogFile", "CacheExpires"}

	for _, name := range stringOpts {
		if viper.IsSet(name) {
			opts[name] = viper.GetString(name)
		}
	}

	// Add ints
	var intOpts = []string{"DocumentConcurrencyLimit", "HTTPConcurrencyLimit", "LogLevel", "ExternalTimeout"}

	for _, name := range intOpts {
		if viper.IsSet(name) {
			opts[name] = viper.GetInt(name)
		}
	}

	// Add string slices
	var stringSliceOpts = []string{"IgnoreURLs", "IgnoreDirs", "StripQueryExcludes"}

	for _, name := range stringSliceOpts {
		if viper.IsSet(name) {
			opts[name] = viper.GetStringSlice(name)
		}
	}

	// Add string maps
	var stringMapOpts  = []string{"HTTPHeaders"}

	for _, name := range stringMapOpts {
		if viper.IsSet(name) {
			opts[name] = viper.GetStringMapString(name)
		}
	}

	return opts
}

func dumpConf(opts optsMap) {
	// Prettily print the provided opts map

	// Capture the maximum key length and a list of all keys
	var maxKeyLength int = 0
	var keys []string
	for key := range opts {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
		keys = append(keys, key)
	}

	// Sort our collected list of keys so we can output them in alpha order
	sort.Strings(keys)

	// Iterate over keys and pretty print
	for _, key := range keys {
		var paddingLength = maxKeyLength - len(key)
		fmt.Println(key, strings.Repeat(".", paddingLength), opts[key])
	}
}
