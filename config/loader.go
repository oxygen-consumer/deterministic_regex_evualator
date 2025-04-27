package config

import (
	"encoding/json"
	"os"
)

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
