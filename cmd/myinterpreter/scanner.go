package main

import (
	"fmt"
	"os"
	"errors"
)

type Scanner struct {
	current int
	contents string
	tokens []Token
}

func (scanner *Scanner) AddToken(token_type TokenType, lexeme string, literal interface{}) {
	var new_token Token
	new_token.token_type = token_type
	new_token.lexeme = lexeme
	new_token.literal = nil
	scanner.tokens = append(scanner.tokens, new_token)
}

func (scanner *Scanner) Advance() byte {
	ch := scanner.contents[scanner.current]
	scanner.current++
	return ch
}

func (scanner *Scanner) Match(char byte) bool {
	if scanner.AtEnd() || scanner.contents[scanner.current] != char {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) AtEnd() bool {
	return scanner.current == len(scanner.contents)
}

func (scanner *Scanner) Scan(lox_file_contents string) error {
	scanner.contents = lox_file_contents
	scanner.current = 0
	line := 1
	var err error
	for ; !scanner.AtEnd(); {
		char := scanner.Advance()
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
		case '\n': line++; break;
		case '=':
			if scanner.Match('=') {
				scanner.AddToken(EQUAL_EQUAL, "==", nil)
			} else {
				scanner.AddToken(EQUAL, "=", nil)
			}
			break;
		case '!':
			if scanner.Match('=') {
				scanner.AddToken(BANG_EQUAL, "!=", nil)
			} else {
				scanner.AddToken(BANG, "!", nil)
			}
			break;
		case '<':
			if scanner.Match('=') {
				scanner.AddToken(LESS_EQUAL, "<=", nil)
			} else {
				scanner.AddToken(LESS, "<", nil)
			}
			break;
		case '>':
			if scanner.Match('=') {
				scanner.AddToken(GREATER_EQUAL, ">=", nil)
			} else {
				scanner.AddToken(GREATER, ">", nil)
			}
			break;
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