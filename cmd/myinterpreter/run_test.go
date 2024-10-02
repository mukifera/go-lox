package main

import (
	"os"

	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
)

type test_config struct {
	Name           string `yaml:"name"`
	FileContents   string `yaml:"fileContents"`
	ExpectedOutput string `yaml:"expectedOutput"`
	ExpectedError  string `yaml:"expectedError"`
}

func TestRuntime(t *testing.T) {

	var tests []test_config

	yamlFile, err := os.ReadFile("run_tests.yaml")
	if err != nil {
		t.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tests)
	if err != nil {
		t.Errorf("Unmarshal: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			exprs := getExpressions(tt.FileContents, t)
			var buf bytes.Buffer
			err := FRunExpressions(&buf, exprs)
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
