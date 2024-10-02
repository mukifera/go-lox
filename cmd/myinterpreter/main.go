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

func handleTokenize() {
	scanner := setupScanner()
	err := scanner.Scan()
	fmt.Println(scanner.StringifyTokens())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(65)
	}
}

func handleParse() {
	parser := setupParser()
	err := parser.Parse()
	if err != nil {
		os.Exit(65)
	}
	fmt.Print(parser.StringifyExpressions())
}

func handleEvaluate() {
	parser := setupParser()
	err := parser.Parse()
	if err != nil {
		os.Exit(65)
	}

	values, err := EvaluateExpressions(parser.expressions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(70)
	}
	for i := 0; i < len(values); i++ {
		value := values[i]
		str := StringifyEvaluationValue(value)
		fmt.Println(str)
	}

}

func handleRun() {
	parser := setupParser()
	err := parser.Parse()
	if err != nil {
		os.Exit(65)
	}

	err = RunExpressions(parser.expressions)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	case "tokenize":
		handleTokenize()
	case "parse":
		handleParse()
	case "evaluate":
		handleEvaluate()
	case "run":
		handleRun()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

}
