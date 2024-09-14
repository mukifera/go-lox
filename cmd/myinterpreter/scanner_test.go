package main

import "testing"

func TestTokenization(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected string
	}{
		{"Empty", "", "EOF  null"},
		{"Parentheses", "(()", "LEFT_PAREN ( null\nLEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
		{"Braces", "{{}}", "LEFT_BRACE { null\nLEFT_BRACE { null\nRIGHT_BRACE } null\nRIGHT_BRACE } null\nEOF  null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.fileContents)

			err := scanner.Scan()
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