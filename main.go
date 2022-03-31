package main

import (
	"fmt"
	"glimmer/executor"
	"os"

	"github.com/pborman/getopt/v2"
)

func main() {
	evalFlag := getopt.BoolLong("repl", 'e', "launch the Glimmer ReadEvalPrintLoop (REPL)")
	parseFlag := getopt.BoolLong("rppl", 'p', "launch the Glimmer ReadParsePrintLoop (RPPL)")
	lexFlag := getopt.BoolLong("rlpl", 'l', "launch the Glimmer ReadLexPrintLoop (RLPL)")
	dotFlag := getopt.BoolLong("dot", 'd', "save the parsed Abstract Syntax Tree as a dotfile and image (infile, repl, and rppl only)")
	outFlag := getopt.BoolLong("output", 'o', "print the evaluated object of the last statement (file option only)")
	getopt.Parse()
	positionalArgs := getopt.Args()

	if moreThanOneServiceSelected(evalFlag, parseFlag, lexFlag) {
		fmt.Println("Error: only one service must be selected")
		printUsageAndDie()
	}

	if *evalFlag || (len(positionalArgs) == 0 && !*parseFlag && !*lexFlag) {
		printService("REPL")
		executor.StartREPL(os.Stdin, os.Stdout, *dotFlag)
	} else if *parseFlag {
		printService("RPPL")
		executor.StartRPPL(os.Stdin, os.Stdout, *dotFlag)
	} else if *lexFlag {
		printService("RLPL")
		executor.StartRLPL(os.Stdin, os.Stdout)
	} else if len(positionalArgs) == 1 {
		evaluated, errs := executor.RunFile(positionalArgs[0], *dotFlag)
		if *outFlag && evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}
		printErrors(errs)
	} else {
		fmt.Println("Error: Only one positional argument <in file> must be given")
		printUsageAndDie()
	}
}

func printService(service string) {
	fmt.Printf("Glimmer %s \u00A9 Daniel Johnson 2022. \"exit\" to exit.\n", service)
}

func printUsageAndDie() {
	getopt.Usage()
	os.Exit(1)
}

func printErrors(errors []error) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func moreThanOneServiceSelected(services ...*bool) bool {
	numSelected := 0
	for _, service := range services {
		if *service {
			numSelected++
		}
	}
	return numSelected > 1
}
