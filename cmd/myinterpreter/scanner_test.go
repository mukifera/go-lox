package main

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestTokenization(t *testing.T) {
	var tests []test_config

	yamlFile, err := os.ReadFile("scanner_tests.yaml")
	if err != nil {
		t.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tests)
	if err != nil {
		t.Errorf("Unmarshal: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			scanner := NewScanner(tt.FileContents)

			err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner: tokenizing error: %v", err)
			}

			actual := scanner.StringifyTokens()
			if actual != tt.ExpectedOutput {
				t.Errorf("Tokenization result is incorrect\nExpected:\n\n%s\n\nGot:\n\n%s", tt.ExpectedOutput, actual)
			}
		})
	}
}
