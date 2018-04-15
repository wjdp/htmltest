package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wjdp/htmltest/output"
	"github.com/wjdp/htmltest/issues"
	"strings"
	"time"
	"github.com/wjdp/htmltest/htmltest"
	"github.com/fatih/color"
	"os"
	"path"
)

const cmdSeparator string = "========================================================================"

var (
	Version   string
	BuildDate string

	confFile    string
	confFileSet bool

	fileMode bool

	dumpVersion bool
	dumpConfig  bool

	logLevel    int
	logLevelSet bool

	skipExternal bool
)

var rootCmd = &cobra.Command{
	Use:   "htmltest PATH",
	Short: "Test generated HTML for problems",
	Long: `htmltest: Test generated HTML for problems. Runs the full suite on PATH.
          Optionally configure everything in a config file instead.
          https://github.com/wjdp/htmltest`,
	Args: cobra.MaximumNArgs(1),
	Run:  runRoot,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&confFile, "conf", "c", ".htmltest.yml",
		"config file")
	rootCmd.PersistentFlags().BoolVarP(&dumpVersion, "version", "v", false,
		"print version and build time")
	rootCmd.PersistentFlags().IntVarP(&logLevel, "log-level", "l", 0,
		"logging level, 0–3: debug, info, warning, error")
	rootCmd.PersistentFlags().BoolVarP(&skipExternal, "skip-external", "s", false,
		"skip external link checks, may shorten execution time considerably.")
	rootCmd.PersistentFlags().BoolVarP(&dumpConfig, "dump-conf", "d", false,
		"print config and exit")

}

func initConfig() {
	// Did the user set the log level? Cannot check with the flag variable as zero is a valid value.
	confFileSet = rootCmd.Flags().Changed("conf")
	logLevelSet = rootCmd.Flags().Changed("log-level")
}

func runRoot(cmd *cobra.Command, args []string) {
	if dumpVersion {
		fmt.Println("htmltest", Version)
		fmt.Println(BuildDate)
		return
	}

	confExists, err := exists(confFile)
	output.CheckErrorPanic(err)
	if confExists {
		// Read config to global viper object
		readConfig(confFile)
	} else {
		if confFileSet {
			output.AbortWith("Cannot open config file '" + confFile + "', no such file")
		}
	}

	if logLevelSet {
		// Override the LogLevel if the flag has been set
		if issues.MinLevel <= logLevel && logLevel <= issues.MaxLevel {
			viper.Set("LogLevel", logLevel)
		} else {
			output.AbortWith("LogLevel must be be between", issues.MinLevel, "and", issues.MaxLevel,
				", provided value is", logLevel)
		}
	}

	// Manually set dir or file to test
	if len(args) == 1 {
		// Open file or directory
		f, err := os.Open(path.Clean(args[0]))
		if os.IsNotExist(err) {
			output.AbortWith("Cannot access '" + args[0] + "', no such file or directory.")
		}
		output.CheckErrorGeneric(err)
		defer f.Close()

		// Get FileInfo, (scan for details)
		fi, err := f.Stat()
		output.CheckErrorPanic(err)

		if fi.IsDir() {
			// We have a directory
			viper.Set("DirectoryPath", path.Clean(args[0]))
		} else {
			// We have a file
			viper.Set("DirectoryPath", path.Dir(args[0]))
			viper.Set("FilePath", path.Base(args[0]))
		}
	}

	// Are we running in file or directory mode?
	fileMode = viper.IsSet("FilePath")

	if skipExternal {
		output.Warn("Skipping the checking of external links.")
		viper.Set("CheckExternal", false)
	}

	// Turn the viper config state into a options map, changes to the viper object after here will have no effect
	opts := viperToOpts()

	if dumpConfig {
		fmt.Println("Dumping config")
		dumpConf(opts)
		return
	}

	// Pass version into options
	opts["Version"] = strings.TrimLeft(Version, "v") // TODO trimming should be removed

	exitCode := runHtmltest(opts)
	os.Exit(exitCode)
}

func runHtmltest(options optsMap) int {
	timeStart := time.Now()

	fmt.Println("htmltest started at", timeStart.Format("03:04:05"), "on", options["DirectoryPath"])
	fmt.Println(cmdSeparator)

	// Run htmltest
	hT, err := htmltest.Test(options)

	if err != nil {
		// Couldn't even run, dump error and exit
		output.AbortWith(err)
	}

	timeEnd := time.Now()
	numErrors := hT.CountErrors()

	if numErrors == 0 {
		color.Set(color.FgHiGreen)
		fmt.Println("✔✔✔ passed in", timeEnd.Sub(timeStart))
		if !fileMode {
			fmt.Println("tested", hT.CountDocuments(), "documents")
		}
		color.Unset()
		return 0
	}

	color.Set(color.FgHiRed)
	fmt.Println(cmdSeparator)
	fmt.Println("✘✘✘ failed in", timeEnd.Sub(timeStart))
	if fileMode {
		fmt.Println(numErrors, "errors")
	} else {
		fmt.Println(numErrors, "errors in", hT.CountDocuments(), "documents")
	}
	color.Unset()
	return 1
}

func Execute() {
	rootCmd.Execute()
}
