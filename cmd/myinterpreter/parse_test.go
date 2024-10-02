package main

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestExpressionParsing(t *testing.T) {
	var tests []test_config

	yamlFile, err := os.ReadFile("parse_tests.yaml")
	if err != nil {
		t.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tests)
	if err != nil {
		t.Errorf("Unmarshal: %v", err)
	}

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
				t.Errorf("Parser: Error while parsing expressions")
			}

			actual := strings.Trim(parser.StringifyExpressions(), "\n")
			if actual != tt.ExpectedOutput {
				t.Errorf("Expression parsing result is incorrect\nExpected:\n\n%s\nGot:\n\n%s\n", tt.ExpectedOutput, actual)
			}
		})
	}
}
