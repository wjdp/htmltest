package main

import(
  "os"
  // "log"
  "htmltest"
  // "issues"
)

func main() {
  bPath := os.Args[1]
  htmltest.SetBasePath(bPath)
  htmltest.Go()
}
