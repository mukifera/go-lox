package main

import (
	"fmt"
	"os"
)

func handleTokenize() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(string(fileContents))
	err = scanner.Scan()
	fmt.Println(scanner.StringifyTokens())
	if err != nil {
		os.Exit(65)
	}
}

func handleParse() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	scanner := NewScanner(string(fileContents))
	err = scanner.Scan()
	if err != nil {
		os.Exit(65)
	}

	parser := NewParser(scanner.tokens)
	parser.Parse()
	fmt.Println(parser.StringifyExpressions())
}

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize": handleTokenize(); break;
	case "parse": handleParse(); break;
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	
}
