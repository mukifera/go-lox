package main

import "testing"

func TestTokenization(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		expected     string
	}{
		{"Empty", "", "EOF  null"},
		{"Parentheses", "(()", "LEFT_PAREN ( null\nLEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
		{"Braces", "{{}}", "LEFT_BRACE { null\nLEFT_BRACE { null\nRIGHT_BRACE } null\nRIGHT_BRACE } null\nEOF  null"},
		{"Single Character Tokens", "({*.,+*});", "LEFT_PAREN ( null\nLEFT_BRACE { null\nSTAR * null\nDOT . null\nCOMMA , null\nPLUS + null\nSTAR * null\nRIGHT_BRACE } null\nRIGHT_PAREN ) null\nSEMICOLON ; null\nEOF  null"},
		{"Assignment And Equality", "={===}", "EQUAL = null\nLEFT_BRACE { null\nEQUAL_EQUAL == null\nEQUAL = null\nRIGHT_BRACE } null\nEOF  null"},
		{"Negation And Inequality", "!!===", "BANG ! null\nBANG_EQUAL != null\nEQUAL_EQUAL == null\nEOF  null"},
		{"Relational Operators", "<<=>>=", "LESS < null\nLESS_EQUAL <= null\nGREATER > null\nGREATER_EQUAL >= null\nEOF  null"},
		{"Comments", "() // Comment", "LEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
		{"Division Operator", "/()", "SLASH / null\nLEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
		{"Whitespaces", "(\t )", "LEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null"},
		{"String Literals", "\"foo baz\"", "STRING \"foo baz\" foo baz\nEOF  null"},
		{"Number Literals", "42 1234.1234", "NUMBER 42 42.0\nNUMBER 1234.1234 1234.1234\nEOF  null"},
		{"Identifiers", "foo bar _hello", "IDENTIFIER foo null\nIDENTIFIER bar null\nIDENTIFIER _hello null\nEOF  null"},
		{"Reserved Words", "and class else false for fun if nil or print return super this true var while", "AND and null\nCLASS class null\nELSE else null\nFALSE false null\nFOR for null\nFUN fun null\nIF if null\nNIL nil null\nOR or null\nPRINT print null\nRETURN return null\nSUPER super null\nTHIS this null\nTRUE true null\nVAR var null\nWHILE while null\nEOF  null"},
		{"Non-ASCII characters", "\"Señor\"", "STRING \"Señor\" Señor\nEOF  null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.fileContents)

			err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner: tokenizing error: %v", err)
			}

			actual := scanner.StringifyTokens()
			if actual != tt.expected {
				t.Errorf("Tokenization result is incorrect\nExpected:\n\n%s\n\nGot:\n\n%s", tt.expected, actual)
			}
		})
	}
}
