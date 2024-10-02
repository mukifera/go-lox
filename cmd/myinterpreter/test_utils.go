package main

import (
	"errors"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type test_config struct {
	Name           string `yaml:"name"`
	FileContents   string `yaml:"fileContents"`
	ExpectedOutput string `yaml:"expectedOutput"`
	ExpectedError  string `yaml:"expectedError"`
}

func fetchYAMLFile(file_name string, t *testing.T) []test_config {
	var tests []test_config
	yamlFile, err := os.ReadFile(file_name)
	if err != nil {
		t.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tests)
	if err != nil {
		t.Errorf("Unmarshal: %v", err)
	}
	return tests
}

func getExpressions(fileContents string, t *testing.T) []Expression {
	scanner := NewScanner(fileContents)
	err := scanner.Scan()
	if err != nil {
		t.Errorf("Scanner: tokenizing error: %v", err)
	}

	parser := NewParser(scanner.tokens)
	err = errors.Join(err, parser.Parse())
	if err != nil {
		t.Errorf("Parser: parsing error: %v", err)
	}

	return parser.expressions
}
