package config

import (
	"deterministic_regex_evaluator/dfa"
	"fmt"
)

func RunTest(test RegexTest, dfaa *dfa.DFA) {
	ok := true

	for _, testString := range test.TestStrings {
		if dfa.RunDFA(dfaa, testString.Input) != testString.Expected {
			ok = false
			fmt.Printf("\tFAIL: %s\n", testString.Input)
		}
	}

	if ok {
		fmt.Printf("\tPASS\n")
	}
}
