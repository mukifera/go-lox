package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRuntime(t *testing.T) {

	tests := fetchYAMLFile("run_tests.yaml", t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			exprs, err := getExpressions(tt.FileContents, t)
			var buf bytes.Buffer
			sub_err := FRunExpressions(&buf, exprs)
			err = errors.Join(err, sub_err)
			actual_output := buf.String()
			actual_err := ""
			if err != nil {
				actual_err = err.Error()
			}

			if tt.ExpectedOutput != actual_output {
				t.Errorf("Execution result mismatch\nExpected output:\n\n%s\nGot:\n\n%s", tt.ExpectedOutput, actual_output)
			}
			if tt.ExpectedError != actual_err {
				t.Errorf("Execution result mismatch\nExpected error:\n\n%s\nGot:\n\n%s", tt.ExpectedError, actual_err)
			}
		})
	}
}
