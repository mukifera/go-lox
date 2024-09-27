package main

import (
	"errors"
	"fmt"
	"os"
)

type Scanner struct {
	current        int
	start          int
	contents       []rune
	tokens         []Token
	reserved_words map[string]TokenType
}

func NewScanner(contents string) *Scanner {
	var scanner Scanner
	scanner.contents = []rune(contents)
	scanner.current = 0
	scanner.start = 0

	scanner.reserved_words = make(map[string]TokenType)
	{
		scanner.reserved_words["and"] = AND
		scanner.reserved_words["class"] = CLASS
		scanner.reserved_words["else"] = ELSE
		scanner.reserved_words["false"] = FALSE
		scanner.reserved_words["for"] = FOR
		scanner.reserved_words["fun"] = FUN
		scanner.reserved_words["if"] = IF
		scanner.reserved_words["nil"] = NIL
		scanner.reserved_words["or"] = OR
		scanner.reserved_words["print"] = PRINT
		scanner.reserved_words["return"] = RETURN
		scanner.reserved_words["super"] = SUPER
		scanner.reserved_words["this"] = THIS
		scanner.reserved_words["true"] = TRUE
		scanner.reserved_words["var"] = VAR
		scanner.reserved_words["while"] = WHILE
	}

	return &scanner
}

func (scanner *Scanner) AddToken(token_type TokenType, lexeme string, literal interface{}) {
	var new_token Token
	new_token.token_type = token_type
	new_token.lexeme = lexeme
	new_token.literal = literal
	scanner.tokens = append(scanner.tokens, new_token)
}

func (scanner *Scanner) Advance() rune {
	ch := scanner.contents[scanner.current]
	scanner.current++
	return ch
}

func (scanner *Scanner) Match(char rune) bool {
	if scanner.AtEnd() || scanner.contents[scanner.current] != char {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) AtEnd() bool {
	return scanner.current == len(scanner.contents)
}

func (scanner *Scanner) Peek() rune {
	if scanner.AtEnd() {
		return 0
	}
	return scanner.contents[scanner.current]
}

func (scanner *Scanner) PeekNext() rune {
	if scanner.current+1 >= len(scanner.contents) {
		return 0
	}
	return scanner.contents[scanner.current+1]
}

func (scanner *Scanner) ScanNumber() {
	for {
		if !isDigit(scanner.Peek()) {
			break
		}
		scanner.Advance()
	}
	if scanner.Peek() == '.' && isDigit(scanner.PeekNext()) {
		scanner.Advance()
		for {
			if !isDigit(scanner.Peek()) {
				break
			}
			scanner.Advance()
		}
	}
	lexeme := scanner.CurrentLexeme()
	scanner.AddToken(NUMBER, lexeme, stringToBigFloat(lexeme))
}

func (scanner *Scanner) ScanIdentifier() {
	for {
		peek := scanner.Peek()
		if !isAlphaNumeric(peek) {
			break
		}
		scanner.Advance()
	}
	lexeme := scanner.CurrentLexeme()
	if token_type, ok := scanner.reserved_words[lexeme]; ok {
		scanner.AddToken(token_type, lexeme, nil)
	} else {
		scanner.AddToken(IDENTIFIER, lexeme, nil)
	}
}

func (scanner *Scanner) CurrentLexeme() string {
	return string(scanner.contents[scanner.start:scanner.current])
}

func (scanner *Scanner) Scan() error {
	line := 1
	found_error := false
	for !scanner.AtEnd() {
		scanner.start = scanner.current
		char := scanner.Advance()
		switch char {
		case '(':
			scanner.AddToken(LEFT_PAREN, "(", nil)
		case ')':
			scanner.AddToken(RIGHT_PAREN, ")", nil)
		case '{':
			scanner.AddToken(LEFT_BRACE, "{", nil)
		case '}':
			scanner.AddToken(RIGHT_BRACE, "}", nil)
		case ',':
			scanner.AddToken(COMMA, ",", nil)
		case '.':
			scanner.AddToken(DOT, ".", nil)
		case '-':
			scanner.AddToken(MINUS, "-", nil)
		case '+':
			scanner.AddToken(PLUS, "+", nil)
		case ';':
			scanner.AddToken(SEMICOLON, ";", nil)
		case '*':
			scanner.AddToken(STAR, "*", nil)
		case '\n':
			line++
		case '=':
			if scanner.Match('=') {
				scanner.AddToken(EQUAL_EQUAL, "==", nil)
			} else {
				scanner.AddToken(EQUAL, "=", nil)
			}
		case '!':
			if scanner.Match('=') {
				scanner.AddToken(BANG_EQUAL, "!=", nil)
			} else {
				scanner.AddToken(BANG, "!", nil)
			}
		case '<':
			if scanner.Match('=') {
				scanner.AddToken(LESS_EQUAL, "<=", nil)
			} else {
				scanner.AddToken(LESS, "<", nil)
			}
		case '>':
			if scanner.Match('=') {
				scanner.AddToken(GREATER_EQUAL, ">=", nil)
			} else {
				scanner.AddToken(GREATER, ">", nil)
			}
		case '/':
			if scanner.Peek() != '/' {
				scanner.AddToken(SLASH, "/", nil)
				break
			}
			for scanner.Peek() != '\n' && scanner.Peek() != 0 {
				scanner.Advance()
			}
		case '"':
			string_literal := ""
			for {
				if scanner.AtEnd() {
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
					found_error = true
					break
				}
				char := scanner.Advance()
				if char == '"' {
					scanner.AddToken(STRING, "\""+string_literal+"\"", string_literal)
					break
				}
				string_literal += string(char)
			}

		case '\t':
		case ' ':
		default:
			if isDigit(char) {
				scanner.ScanNumber()
			} else if isAlpha(char) {
				scanner.ScanIdentifier()
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, char)
				found_error = true
			}
		}
	}
	scanner.AddToken(EOF, "", nil)
	if found_error {
		return errors.New("error scanning file contents")
	}
	return nil
}

func (scanner Scanner) StringifyTokens() string {
	ret := ""
	for _, token := range scanner.tokens {
		if len(ret) > 0 {
			ret += "\n"
		}
		ret += token.String()
	}
	return ret
}
