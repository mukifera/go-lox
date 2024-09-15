package main

import "testing"

func TestEvaluation(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected []string
	}{
		{"Literals: Boolean/true", "true", []string{"true"}},
		{"Literals: Boolean/false", "false", []string{"false"}},
		{"Literals: nil", "nil", []string{"nil"}},
		{"Literals: Number #1", "10.40", []string{"10.4"}},
		{"Literals: Number #2", "10", []string{"10"}},
		{"Parentheses #1", `("hello world!")`, []string{"hello world!"}},
		{"Parentheses #2", "(true)", []string{"true"}},
		{"Parentheses #3", "(10.40)", []string{"10.4"}},
		{"Parentheses #4", "((false))", []string{"false"}},
		{"Unary: Negation", "-73", []string{"-73"}},
		{"Unary: Not #1", "!true", []string{"false"}},
		{"Unary: Not #2", "!10.40", []string{"false"}},
		{"Unary: Not #3", "!((false))", []string{"true"}},
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
				t.Errorf("Parser: parsing error: %v", err)
			}

			evaluator := NewEvaluator(parser.expressions)
			evaluator.Evaluate()

			actual := evaluator.StringifyValues()

			if len(tt.expected) != len(actual) {
				t.Errorf("Evaluation result length mismatch: Expected %d outputs, got %d", len(tt.expected), len(actual))
			} else {
				for index, str := range actual {
					if tt.expected[index] != str {
						t.Errorf("Evaluation result mismatch on output %d\nExpected: %s\nGot: %s", index + 1, tt.expected[index], str)
					}
				}
			}
		})
	}
}