package config

import "fmt"

func RunTests(tests []RegexTest) error {
	for _, test := range tests {
		// TODO:
		// load regex to dfa
		// run dfa for each string and check the expected
		// if the expected is different than the dfa result, return an error
		fmt.Println(test)
	}

	return nil
}
