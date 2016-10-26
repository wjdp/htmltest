package test

import (
  "strings"
  "github.com/imdario/mergo"
  "github.com/wjdp/htmltest/issues"
)

type Options struct {
  DirectoryPath string
  FilePath string

  CheckExternal bool
  CheckInternal bool
  CheckMailto bool
  CheckTel bool
  EnforceHTTPS bool

  TestFilesConcurrently bool
  LogLevel int

  DirectoryIndex string

  ExternalTimeout int
  StripQueryString bool
  StripQueryExcludes []string
}

func DefaultOptions() Options {
  // Specify defaults here
  options := Options{
    CheckExternal: true,
    CheckInternal: true,
    CheckMailto: true,
    CheckTel: true,
    EnforceHTTPS: false,

    TestFilesConcurrently: false,
    LogLevel: issues.INFO,

    DirectoryIndex: "index.html",

    ExternalTimeout: 1,
    StripQueryString: true,
    StripQueryExcludes: []string{"fonts.googleapis.com"},
  }
  return options
}

func OptionsSetDefaults(opts *Options) {
  mergo.Merge(opts, DefaultOptions())
}

func InList(list []string, key string) bool {
  for _, item := range list {
    if strings.Contains(key, item) { return true }
  }
  return false
}
