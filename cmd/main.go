package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/wjdp/htmltest/htmltest"
	"os"
	"time"
)

const VERSION string = "0.0.1"
const SEPERATOR string = "========================================================================"

func main() {
	usage := `htmltest - Test generated HTML for problems
           https://github.com/wjdp/htmltest

Usage:
  htmltest [--log-level=LEVEL] [<path>]
  htmltest --conf=CFILE
  htmltest --version
  htmltest -h --help

Options:
  <path>              Path to directory or file to test, omit for current
                      directory.
  -h --help           Show this text.
  --log-level=LEVEL   Logging level, 0-3: debug, info, warning, error.
  --conf=CFILE        Use a JSON configuration file for advanced options.`
	versionText := "htmlproofer " + VERSION
	arguments, _ := docopt.Parse(usage, nil, true, versionText, false)
	// fmt.Println(arguments)

	var options map[string]interface{}
	if arguments["--conf"] != nil {
		options = parseConfFile(arguments["--conf"].(string))
	} else {
		options = parseCLIArgs(arguments)
	}

	exitCode := run(options)
	os.Exit(exitCode)

}

type optsMap map[string]interface{}

func parseConfFile(path string) optsMap {
	// Read in json config file
	f, err := os.Open(path)
	checkErr(err)
	defer f.Close()

	var optsJson optsMap
	err = json.NewDecoder(f).Decode(&optsJson)
	checkErr(err)

	return optsJson
}

func parseCLIArgs(arguments map[string]interface{}) optsMap {
	// Deal with cl arguments
	options := map[string]interface{}{}

	if arguments["<path>"] != nil {
		options["DirectoryPath"] = arguments["<path>"]
	} else {
		// All other options exhausted, run on current directory
		options["DirectoryPath"] = "."
	}

	if arguments["--log-level"] != nil {
		options["LogLevel"] = arguments["--log-level"]
	}
	return options
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(options optsMap) int {
	timeStart := time.Now()

	fmt.Println("htmltest started at", timeStart.Format("03:04:05"), "on", options["DirectoryPath"])
	fmt.Println(SEPERATOR)

	hT := htmltest.Test(options)

	timeEnd := time.Now()
	numErrors := hT.CountErrors()

	if numErrors == 0 {
		color.Set(color.FgHiGreen)
		fmt.Println("✔✔✔ passed in", timeEnd.Sub(timeStart))
		color.Unset()
		return 0
	} else {
		color.Set(color.FgHiRed)
		fmt.Println(SEPERATOR)
		fmt.Println("✘✘✘ failed in", timeEnd.Sub(timeStart))
		fmt.Println(numErrors, "errors")
		color.Unset()
		return 1
	}
}
