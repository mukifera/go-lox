package main

import (
	"strings"
	"testing"
)

func TestExpressionParsing(t *testing.T) {
	tests := fetchYAMLFile("parse_tests.yaml", t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			scanner := NewScanner(tt.FileContents)

			err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner: tokenizing error: %v", err)
			}

			parser := NewParser(scanner.tokens)

			err = parser.Parse()
			if err != nil {
				actual_err := err.Error()
				if actual_err != tt.ExpectedError {
					t.Errorf("Expression parsing mismatch\nExpected error:\n\n%s\nGot:\n\n%s", tt.ExpectedError, actual_err)
				}
			} else {
				actual := strings.Trim(parser.StringifyExpressions(), "\n")
				if actual != tt.ExpectedOutput {
					t.Errorf("Expression parsing result is incorrect\nExpected:\n\n%s\nGot:\n\n%s\n", tt.ExpectedOutput, actual)
				}
			}

		})
	}
}
