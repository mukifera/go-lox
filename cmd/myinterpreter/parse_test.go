package main

import "testing"

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected string
	}{
		{"Booleans/true", "true", "true"},
		{"Booleans/false", "false", "false"},
		{"Nil", "nil", "nil"},
		{"Number Literals", "42.47", "42.47"},
		{"String Literals", "\"hello\"", "hello"},
		{"Parentheses", "(\"foo\")", "(group foo)"},
		{"Unary/Negation", "-5", "(- 5.0)"},
		{"Unary/Not", "!true", "(! true)"},
		{"Arithmetic/Multiplication And Division", "16 * 38 / 58", "(/ (* 16.0 38.0) 58.0)"},
		{"Arithmetic/Addtion And Subtraction", "52 + 80 - 94", "(- (+ 52.0 80.0) 94.0)"},
		{"Comparision Operators", "83 < 99 > 115 <= 11 >= 1", "(>= (<= (> (< 83.0 99.0) 115.0) 11.0) 1.0)"},
		{"Equality Operators", `"baz" == "baz" != "bar"`, "(!= (== baz baz) bar)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.fileContents)

			err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner: tokenizing error: %v", err)
			}

			parser := NewParser(scanner.tokens)

			err = parser.Parse()
			if err != nil {
				t.Errorf("Parser: Error while parsing expressions")
			}

			actual := parser.StringifyExpressions()
			if actual != tt.expected {
				t.Errorf("Expression parsing result is incorrect\nExpected:\n\n%s\n\nGot:\n\n%s", tt.expected, actual)
			}
		})
	}
}