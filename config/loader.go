package config

import (
	"encoding/json"
	"os"
)

type TestString struct {
	Input    string `json:"input"`
	Expected bool   `json:"expected"`
}

type RegexTest struct {
	Name        string       `json:"name"`
	Regex       string       `json:"regex"`
	TestStrings []TestString `json:"test_strings"`
}

func LoadTests(filename string) ([]RegexTest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tests []RegexTest
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&tests); err != nil {
		return nil, err
	}

	return tests, nil
}
