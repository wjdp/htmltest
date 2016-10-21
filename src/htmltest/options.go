package htmltest

import "issues"

type Options struct {
  CheckExternal bool
  CheckInternal bool
  CheckMailto bool
  CheckTel bool
  EnforceHTTPS bool

  TestFilesConcurrently bool
  LogLevel int

  DirectoryIndex string
}

func NewOptions() Options {
  // Specify defaults here
  options := Options{
    CheckExternal: false,
    CheckInternal: true,
    CheckMailto: true,
    CheckTel: true,
    EnforceHTTPS: false,

    TestFilesConcurrently: false,
    LogLevel: issues.INFO,

    DirectoryIndex: "index.html",
  }
  return options
}
