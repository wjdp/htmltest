package main

import (
	"github.com/wjdp/htmltest/htmltest"
	"log"
	"os"
	// "issues"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid argument")
	}

	options := map[string]interface{}{
		"DirectoryPath": os.Args[1],
	}

	htmltest.Test(options)
}
