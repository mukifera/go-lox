package main

import (
	"fmt"
	"os"
	"errors"
)

type Scanner struct {
	tokens []Token
}

func (scanner *Scanner) AddToken(token_type TokenType, lexeme string, literal interface{}) {
	var new_token Token
	new_token.token_type = token_type
	new_token.lexeme = lexeme
	new_token.literal = nil
	scanner.tokens = append(scanner.tokens, new_token)
}

func (scanner *Scanner) Scan(lox_file_contents string) error {
	line := 1
	var err error
	err = nil
	for i := 0; i < len(lox_file_contents); i++ {
		char := lox_file_contents[i]
		switch char {
		case '(': scanner.AddToken(LEFT_PAREN, "(", nil); break;
		case ')': scanner.AddToken(RIGHT_PAREN, ")", nil); break;
		case '{': scanner.AddToken(LEFT_BRACE, "{", nil); break;
		case '}': scanner.AddToken(RIGHT_BRACE, "}", nil); break;
		case ',': scanner.AddToken(COMMA, ",", nil); break;
		case '.': scanner.AddToken(DOT, ".", nil); break;
		case '-': scanner.AddToken(MINUS, "-", nil); break;
		case '+': scanner.AddToken(PLUS, "+", nil); break;
		case ';': scanner.AddToken(SEMICOLON, ";", nil); break;
		case '*': scanner.AddToken(STAR, "*", nil); break;
		case '=':
			if i+1 < len(lox_file_contents) && lox_file_contents[i+1] == '='{
				scanner.AddToken(EQUAL_EQUAL, "==", nil)
				i++;
			} else {
				scanner.AddToken(EQUAL, "=", nil)
			}
			break;
		case '\n': line++; break;
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, char)
			err = errors.New("Unexpected characters")
		}
	}
	scanner.AddToken(EOF, "", nil)
	return err
}

func (scanner Scanner) StringifyTokens() string {
	ret := ""
	for _, token := range scanner.tokens {
		if len(ret) > 0 {
			ret += "\n"
		}
		ret += token.String()
	}
	return ret;
}