package main

import (
	"os"

	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
)

type test_config struct {
	Name         string `yaml:"name"`
	FileContents string `yaml:"fileContents"`
	Expected     string `yaml:"expected"`
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
			FRunExpressions(&buf, exprs)
			str := buf.String()

			if tt.Expected != str {
				t.Errorf("Execution result mismatch\nExpected:\n\n%s\nGot:\n\n%s", tt.Expected, str)
			}
		})
	}
}
