package main

import (
	"fmt"
	"os"
)

func readFile() []byte {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return fileContents
}

func setupScanner() *Scanner {
	fileContents := readFile()
	return NewScanner(string(fileContents))
}

func setupParser() *Parser {
	scanner := setupScanner()
	err := scanner.Scan()
	if err != nil {
		os.Exit(65)
	}
	return NewParser(scanner.tokens)
}

func setupEvaluator() *Evaluator {
	parser := setupParser()
	err := parser.Parse()
	if err != nil {
		os.Exit(65)
	}
	return NewEvaluator(parser.expressions)
}

func handleTokenize() {
	scanner := setupScanner()
	err := scanner.Scan()
	fmt.Println(scanner.StringifyTokens())
	if err != nil {
		os.Exit(65)
	}
}

func handleParse() {
	parser := setupParser()
	err := parser.Parse()
	if err != nil {
		os.Exit(65)
	}
	fmt.Println(parser.StringifyExpressions())
}

func handleEvaluate() {
	evaluator := setupEvaluator()
	evaluator.Evaluate()
	strs := evaluator.StringifyValues()
	found_error := false
	for index, str := range strs {
		err := evaluator.errors[index]
		if err != nil {
			found_error = true
			fmt.Println(err)
		} else {
			fmt.Println(str)
		}
	}
	if found_error {
		os.Exit(70)
	}
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
	case "evaluate": handleEvaluate(); break;
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	
}
