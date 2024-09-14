package main

import (
	"os"
	"fmt"
	"testing"
)

func hookupScanner(filename string) *Scanner {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	return NewScanner(string(fileContents))
}

func TestScanEmpty(t *testing.T) {
	scanner := hookupScanner("test_files/empty.lox")
	err := scanner.Scan()
	if err != nil {
		t.Errorf("Error building Scanner")
	}
	actual := scanner.StringifyTokens()
	expected := "EOF  null"
	if actual != expected {
		t.Errorf("Tokenization result is incorrect\nGot: %s\nExpected: %s", actual, expected)
	}
}