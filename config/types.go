package config

type TestString struct {
	Input    string `json:"input"`
	Expected bool   `json:"expected"`
}

type RegexTest struct {
	Name        string       `json:"name"`
	Regex       string       `json:"regex"`
	TestStrings []TestString `json:"test_strings"`
}
