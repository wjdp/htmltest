// htmltest CLI
package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/wjdp/htmltest/htmltest"
	"github.com/wjdp/htmltest/output"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const cmdVersion string = "0.3.1"
const cmdSeparator string = "========================================================================"

var (
	buildDate string
)

func main() {
	usage := `htmltest - Test generated HTML for problems
           https://github.com/wjdp/htmltest

Usage:
  htmltest [options] [<path>]
  htmltest -v --version
  htmltest -h --help

Options:
  <path>                       Path to directory or file to test, if omitted we
                               attempt to read from .htmltest.yml.
  -c FILE, --conf FILE         Custom path to config file.
  -h, --help                   Show this text.
  -l LEVEL, --log-level LEVEL  Logging level, 0-3: debug, info, warning, error.
  -s, --skip-external          Skip external link checks, may shorten execution
                               time considerably.
  -v, --version                Show version and build time.
`
	versionText := "htmltest " + cmdVersion + "\n" + buildDate
	arguments, _ := docopt.Parse(usage, nil, true, versionText, false)

	// fmt.Println(arguments)

	var options map[string]interface{}
	if arguments["--conf"] != nil {
		// Config file specified
		options = parseConfFile(arguments, arguments["--conf"].(string), true)
	} else if arguments["<path>"] != nil {
		// Path specified
		options = parseCLIArgs(arguments)
	} else {
		// Other
		options = parseConfFile(arguments, ".htmltest.yml", false)
	}

	exitCode := run(options)
	os.Exit(exitCode)

}

type optsMap map[string]interface{}

func parseConfFile(arguments map[string]interface{}, path string, explicit bool) optsMap {
	yamlFile, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		if explicit {
			output.AbortWith("The provided config file: ", path, "does not exist.")
		} else {
			output.AbortWith(`No path provided & the default config .htmltest.yml does not exist.
See htmltest -h for usage.`)
		}
	}
	output.CheckErrorGeneric(err)

	var optsConf optsMap
	err = yaml.Unmarshal(yamlFile, &optsConf)
	output.CheckErrorGeneric(err)

	// Override or append config options with any specified in CLI args
	augmentWithCLIArgs(optsConf, arguments)

	return optsConf
}

// Wrapper for augmentWithCLIArgs when you don't already have an options map.
func parseCLIArgs(arguments map[string]interface{}) optsMap {
	options := optsMap{}
	augmentWithCLIArgs(options, arguments)
	return options
}

// Override or append to the options map with CLI args.
func augmentWithCLIArgs(options optsMap, arguments map[string]interface{}) {
	// Deal with cli arguments

	if arguments["<path>"] != nil {
		options["DirectoryPath"] = arguments["<path>"]
	}

	if arguments["--log-level"] != nil {
		if ll, err := strconv.Atoi(arguments["--log-level"].(string)); err == nil && ll >= 0 {
			options["LogLevel"] = ll
		} else {
			output.AbortWith("--log-level must be a positive integer")
		}
	}

	if arguments["--skip-external"].(bool) {
		output.Warn("Skipping the checking of external links.")
		options["CheckExternal"] = false
	}

}

func run(options optsMap) int {
	timeStart := time.Now()

	fmt.Println("htmltest started at", timeStart.Format("03:04:05"), "on", options["DirectoryPath"])
	fmt.Println(cmdSeparator)

	hT := htmltest.Test(options)

	timeEnd := time.Now()
	numErrors := hT.CountErrors()

	if numErrors == 0 {
		color.Set(color.FgHiGreen)
		fmt.Println("✔✔✔ passed in", timeEnd.Sub(timeStart))
		color.Unset()
		return 0
	}

	color.Set(color.FgHiRed)
	fmt.Println(cmdSeparator)
	fmt.Println("✘✘✘ failed in", timeEnd.Sub(timeStart))
	fmt.Println(numErrors, "errors")
	color.Unset()
	return 1

}
