package main

import (
	"bytes"
	"testing"
)

func TestRuntime(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected string
	}{
		{"Print #1", `print "Hello, World!";`, "Hello, World!\n"},
		{"Print #2", "print 42;", "42\n"},
		{"Print #3", "print true;", "true\n"},
		{"Print #4", "print 12 + 24;", "36\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprs := getExpressions(tt.fileContents, t)
			var buf bytes.Buffer
			FRunExpressions(&buf, exprs)
			str := buf.String()

			if tt.expected != str {
				t.Errorf("Execution result mismatch\nExpected:\n\n%s\nGot:\n\n%s", tt.expected, str)
			}
		})
	}
}