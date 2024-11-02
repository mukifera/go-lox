package main

import (
	"strings"
	"testing"
)

func TestEvaluation(t *testing.T) {
	tests := fetchYAMLFile("evaluate_tests.yaml", t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			expr := getExpression(tt.FileContents, t)
			values, err := EvaluateExpressions([]Expression{expr})
			actual := StringifyEvaluationValues(values)

			if err != nil {
				actual_err := err.Error()
				actual_err = strings.Trim(actual_err, "\n")

				expected_err := strings.Trim(tt.ExpectedError, "\n")

				if expected_err != actual_err {
					t.Errorf("Evaluation error mismatch\nExpected: %s\nGot: %s", expected_err, actual_err)
				}
			} else {
				if tt.ExpectedOutput != actual {
					t.Errorf("Evaluation result mismatch\nExpected: %s\nGot: %s", tt.ExpectedOutput, actual)
				}
			}

		})
	}
}
