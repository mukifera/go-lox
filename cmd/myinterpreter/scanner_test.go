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
		{"Single Character Tokens", "({*.,+*});", "LEFT_PAREN ( null\nLEFT_BRACE { null\nSTAR * null\nDOT . null\nCOMMA , null\nPLUS + null\nSTAR * null\nRIGHT_BRACE } null\nRIGHT_PAREN ) null\nSEMICOLON ; null\nEOF  null"},
		{"Assignment And Equality", "={===}", "EQUAL = null\nLEFT_BRACE { null\nEQUAL_EQUAL == null\nEQUAL = null\nRIGHT_BRACE } null\nEOF  null"},
		{"Negation And Inequality", "!!===", "BANG ! null\nBANG_EQUAL != null\nEQUAL_EQUAL == null\nEOF  null"},
		{"Relational Operators", "<<=>>=", "LESS < null\nLESS_EQUAL <= null\nGREATER > null\nGREATER_EQUAL >= null\nEOF  null"},
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