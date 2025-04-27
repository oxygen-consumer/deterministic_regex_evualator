package config

import (
	"deterministic_regex_evaluator/nfa"
	"deterministic_regex_evaluator/regex"
	"fmt"
	"os"
)

func RunTests(tests []RegexTest) error {
	for _, test := range tests {
		// TODO:
		// load regex to dfa
		// run dfa for each string and check the expected
		// if the expected is different than the dfa result, return an error

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

		dotString := builtNFA.ToDot()
		os.WriteFile("./out/"+test.Name+"_nfa.dot", []byte(dotString), 0644)
	}

	return nil
}
