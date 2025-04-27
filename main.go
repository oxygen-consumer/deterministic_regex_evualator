package main

import (
	"deterministic_regex_evaluator/config"
	"flag"
	"fmt"
	"os"
)

func main() {
	var testsFilename string

	flag.CommandLine.StringVar(&testsFilename, "tests", "", "Path to the tests file")
	flag.Parse()

	if testsFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	tests, err := config.LoadTests(testsFilename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	err = config.RunTests(tests)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
