package main

import (
	"os"
	"testing"
)

func TestTokenization(t *testing.T) {
	tests := []struct {
		name string
		filename string
		expected string
	}{
		{"Empty", "empty.lox", "EOF  null"},
		{"Parentheses", "parentheses.lox", "LEFT_PAREN ( null\nLEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fileContents, err := os.ReadFile("test_files/" + tt.filename)
			if err != nil {
				t.Errorf("Error reading file: %v\n", err)
			}
			scanner := NewScanner(string(fileContents))

			err = scanner.Scan()
			if err != nil {
				t.Errorf("Error building Scanner")
			}

			actual := scanner.StringifyTokens()
			if actual != tt.expected {
				t.Errorf("Tokenization result is incorrect\nExpected:\n\n%s\n\nGot:\n\n%s", tt.expected, actual)
			}
		})
	}
}