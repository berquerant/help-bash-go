package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	"github.com/berquerant/logger"
)

var (
	displayFunctions        = flag.Bool("f", false, "Displays top-level (outside of any statements) functions documentations.")
	displayVariables        = flag.Bool("r", false, "Displays top-level variables documentations.")
	displayVersion          = flag.Bool("v", false, "Prints version information.")
	displayCommentThreshold = flag.Uint("t", 3, "Number of rows of which threshold to determine whether to display the top-level comments.\ne.g. the value is 3 then displays 3 or more lines of top-level comments.")
)

// flag's name should be one letter due to fixArgs()

func isBoolFlag(f *flag.Flag) bool {
	// https://pkg.go.dev/flag#Value
	_, found := reflect.TypeOf(f.Value).MethodByName("IsBoolFlag")
	return found
}

func fixArgs(args []string) []string {
	var (
		flagChars        = "h"
		flagWithArgChars string
	)
	flag.VisitAll(func(f *flag.Flag) {
		if isBoolFlag(f) {
			flagChars += f.Name
		} else {
			flagWithArgChars += f.Name
		}
	})
	return FixArgs(args, flagChars, flagWithArgChars)
}

const usage = `help-bash displays the things like documentation comments in .sh file.

USAGE:
  help-bash [OPTIONS] [PATH]


ARGS:
    <PATH>
        A file to display documentations.


ENVIRONMENT VARIABLES:
     HELP_BASH_DEBUG
         If value is 1 then enables debug logs.
         If value is 2 then enables trace logs in addition to debug logs.


OPTIONS:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func initLogger() {
	debugLevel, err := strconv.Atoi(os.Getenv("HELP_BASH_DEBUG"))
	if err != nil {
		return
	}
	if debugLevel > 0 {
		logger.G().SetLevel(logger.Ldebug)
	}
	if debugLevel > 1 {
		logger.G().SetLevel(logger.Ltrace)
	}
}

func getTargetFile() string {
	if flag.NArg() > 0 {
		return flag.Args()[0]
	}
	return ""
}

func debugParams(targetFile string) {
	flag.VisitAll(func(f *flag.Flag) {
		logger.G().Debug("Parameter:%s=%v", f.Name, f.Value)
	})
	logger.G().Debug("Parameter:targetFile:%s", targetFile)
}

func main() {
	initLogger()
	logger.G().Trace("Raw args: %v", os.Args)
	os.Args = fixArgs(os.Args)
	logger.G().Trace("Fixed args: %v", os.Args)
	flag.Usage = Usage
	flag.Parse()
	os.Exit(run())
}

func run() int {
	targetFile := getTargetFile()
	debugParams(targetFile)

	if *displayVersion {
		printVersion()
		return 0
	}

	execExtract := func(w io.Writer, r io.Reader) int {
		if err := Extract(
			w,
			r,
			WithCommentThreshold(*displayCommentThreshold),
			WithNeedFunc(*displayFunctions),
			WithNeedVar(*displayVariables),
		); err != nil {
			logger.G().Error("%v", err)
			return 1
		}
		return 0
	}

	if targetFile == "" {
		return execExtract(os.Stdout, os.Stdin)
	}

	f, err := os.Open(targetFile)
	if err != nil {
		logger.G().Error("%v", err)
		return 1
	}
	defer f.Close()
	return execExtract(os.Stdout, f)
}

func printVersion() {
	fmt.Println(Ldflags())
}
