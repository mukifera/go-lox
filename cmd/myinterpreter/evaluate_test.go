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
		{"Literals: Boolean/nil", "nil", []string{"nil"}},
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