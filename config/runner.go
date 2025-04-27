package config

import (
	"deterministic_regex_evaluator/regex"
	"fmt"
)

func RunTests(tests []RegexTest) error {
	for _, test := range tests {
		// TODO:
		// load regex to dfa
		// run dfa for each string and check the expected
		// if the expected is different than the dfa result, return an error

		tokens := regex.Tokenize(test.Regex)

		postfix, err := regex.ToPostfix(tokens)
		if err != nil {
			return err
		}

		var postfixRunes []rune
		for _, tok := range postfix {
			postfixRunes = append(postfixRunes, tok.Value)
		}

		fmt.Printf("test %s: %s -> %s\n", test.Name, test.Regex, string(postfixRunes))
	}

	return nil
}
