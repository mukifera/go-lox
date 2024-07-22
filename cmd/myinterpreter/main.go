package main

import (
	"fmt"
	"os"
	"strings"
)

func tokenize(lox_file_contents string) []string {
	var tokenization []string
	for i := 0; i < len(lox_file_contents); i++ {
		char := lox_file_contents[i]
		if char == '('{
			tokenization = append(tokenization, "LEFT_PAREN ( null")
		}else if char == ')' {
			tokenization = append(tokenization, "RIGHT_PAREN ) null")
		}else if char == '{' {
			tokenization = append(tokenization, "LEFT_BRACE ) null")
		}else if char == '}' {
			tokenization = append(tokenization, "RIGHT_BRACE ) null")
		}
	}
	tokenization = append(tokenization, "EOF  null")
	return tokenization
}

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println(strings.Join(tokenize(string(fileContents)), "\n"))
}
