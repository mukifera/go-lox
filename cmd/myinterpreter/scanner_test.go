package main

import (
	"testing"
)

func TestTokenization(t *testing.T) {
	tests := fetchYAMLFile("scanner_tests.yaml", t)

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
