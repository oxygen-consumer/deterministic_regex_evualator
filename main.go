package main

import (
	"deterministic_regex_evaluator/config"
	"deterministic_regex_evaluator/dfa"
	"deterministic_regex_evaluator/nfa"
	"deterministic_regex_evaluator/regex"
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

	for _, test := range tests {
		fmt.Printf("Test %s: %s\n", test.Name, test.Regex)

		builtDFA, err := regexToDFAWithSubsetConstruction(test.Regex)
		if err != nil {
			fmt.Printf("\tERROR: %s\n", err)
		}

		config.RunTest(test, builtDFA)
	}
}

func regexToDFAWithSubsetConstruction(regexx string) (*dfa.DFA, error) {
	tokens := regex.Tokenize(regexx)

	postfix, err := regex.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to postfix: %s", err)
	}

	builtNFA, err := nfa.BuildNFA(postfix)
	if err != nil {
		return nil, fmt.Errorf("failed to build nfa: %s", err)
	}

	return dfa.BuildDFAFromNFA(builtNFA), nil
}
