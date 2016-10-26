package main

import(
  "os"
  // "log"
  "github.com/wjdp/htmltest/test"
  // "issues"
)

func main() {
  bPath := os.Args[1]
  test.SetBasePath(bPath)
  test.Setup()
  test.Go()
}
