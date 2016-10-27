package main

import(
  "os"
  "log"
  "github.com/wjdp/htmltest/test"
  // "issues"
)

func main() {
  if len(os.Args) != 2 {
    log.Fatal("Invalid argument")
  }

  options := map[string]interface{}{
    "DirectoryPath": os.Args[1],
  }

  test.Test(options)
}
