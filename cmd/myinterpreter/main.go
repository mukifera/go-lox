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
	errs := scanner.Scan()
	if len(errs) != 0 {
		os.Exit(65)
	}
	return NewParser(scanner.tokens)
}

func handleTokenize() {
	scanner := setupScanner()
	errs := scanner.Scan()
	fmt.Println(scanner.StringifyTokens())
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(65)
	}
}

func handleParse() {
	parser := setupParser()
	errs := parser.Parse()
	if len(errs) != 0 {
		os.Exit(65)
	}
	fmt.Println(parser.StringifyExpressions())
}

func handleEvaluate() {
	parser := setupParser()
	errs := parser.Parse()
	if len(errs) != 0 {
		os.Exit(65)
	}

	values, errs := EvaluateExpressions(parser.expressions)
	found_error := false
	for i := 0; i < len(values); i++ {
		value := values[i]
		err := errs[i]
		str := StringifyEvaluationValue(value)
		if err != nil {
			found_error = true
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			fmt.Println(str)
		}
	}

	if found_error {
		os.Exit(70)
	}
}

func handleRun() {
	parser := setupParser()
	errs := parser.Parse()
	if len(errs) != 0 {
		os.Exit(65)
	}

	errs = RunExpressions(parser.expressions)

	found_error := false
	for _, err := range errs {
		if err != nil {
			found_error = true
			fmt.Fprintf(os.Stderr, "%v\n", err)
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
