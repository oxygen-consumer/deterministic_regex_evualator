package config

import (
	"deterministic_regex_evaluator/dfa"
	"deterministic_regex_evaluator/nfa"
	"deterministic_regex_evaluator/regex"
	"fmt"
	"os"
)

func RunTests(tests []RegexTest) error {
	for _, test := range tests {
		tokens := regex.Tokenize(test.Regex)

		postfix, err := regex.Parse(tokens)
		if err != nil {
			return err
		}

		var postfixRunes []rune
		for _, tok := range postfix {
			postfixRunes = append(postfixRunes, tok.Value)
		}

		fmt.Printf("test %s: %s -> %s\n", test.Name, test.Regex, string(postfixRunes))

		builtNFA, err := nfa.BuildNFA(postfix)
		if err != nil {
			return err
		}

		builtDFA := dfa.BuildDFAFromNFA(builtNFA)
		for _, testString := range test.TestStrings {
			if testString.Expected != dfa.RunDFA(builtDFA, testString.Input) {
				fmt.Printf("test failed for test %s on input %s\n", test.Name, testString.Input)
			}
		}

		dotStringNFA := builtNFA.ToDot()
		os.WriteFile("./out/"+test.Name+"_nfa.dot", []byte(dotStringNFA), 0644)
		dotStringDFA := builtDFA.ToDot()
		os.WriteFile("./out/"+test.Name+"_dfa.dot", []byte(dotStringDFA), 0644)
	}

	return nil
}
