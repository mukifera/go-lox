package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"
)

type Scanner struct {
	current int
	start int
	contents string
	tokens []Token
}

func (scanner *Scanner) AddToken(token_type TokenType, lexeme string, literal interface{}) {
	var new_token Token
	new_token.token_type = token_type
	new_token.lexeme = lexeme
	new_token.literal = literal
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

func (scanner *Scanner) Peek() byte {
	if scanner.AtEnd() {
		return 0
	}
	return scanner.contents[scanner.current]
}

func (scanner *Scanner) PeekNext() byte {
	if scanner.current + 1 >= len(scanner.contents) {
		return 0
	}
	return scanner.contents[scanner.current + 1]
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'z') || (c == '_')
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func stringToBigFloat(str string) big.Float {
	var power, digit, ten, tmp, float_literal big.Float
	power.SetFloat64(1.0)
	ten.SetFloat64(10.0)
	float_literal.SetFloat64(0)
	found_period := false
	for _, char := range str {
		if char == '.' {
			found_period = true
			continue
		}
		if found_period {
			digit.SetFloat64(float64(char - '0'))
			power.Quo(&power, &ten)
			tmp.Mul(&power, &digit)
			float_literal.Add(&float_literal, &tmp)
		} else {
			digit.SetFloat64(float64(char - '0'))
			tmp.Mul(&float_literal, &ten)
			float_literal.Add(&tmp, &digit)
		}
	}
	return float_literal
}

func (scanner *Scanner) ScanNumber() {
	for {
		if !isDigit(scanner.Peek()) { break }
		scanner.Advance()
	}
	if scanner.Peek() == '.' && isDigit(scanner.PeekNext()) {
		scanner.Advance()
		for {
			if !isDigit(scanner.Peek()) { break }
			scanner.Advance()
		}
	}
	lexeme := scanner.contents[scanner.start : scanner.current]
	scanner.AddToken(NUMBER, lexeme, stringToBigFloat(lexeme))
}

func (scanner *Scanner) Scan(lox_file_contents string) error {
	scanner.contents = lox_file_contents
	scanner.current = 0
	line := 1
	found_error := false

	reserved_words := make(map[string]TokenType)
	{
		reserved_words["and"] 	 = AND
		reserved_words["class"]  = CLASS
		reserved_words["else"] 	 = ELSE
		reserved_words["false"]  = FALSE
		reserved_words["for"] 	 = FOR
		reserved_words["fun"] 	 = FUN
		reserved_words["if"] 		 = IF
		reserved_words["nil"] 	 = NIL
		reserved_words["or"] 		 = OR
		reserved_words["print"]  = PRINT
		reserved_words["return"] = RETURN
		reserved_words["super"]  = SUPER
		reserved_words["this"]   = THIS
		reserved_words["true"] 	 = TRUE
		reserved_words["var"] 	 = VAR
		reserved_words["while"]  = WHILE
	}

	for ; !scanner.AtEnd(); {
		scanner.start = scanner.current
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
		case '/':
			if scanner.Peek() != '/' {
				scanner.AddToken(SLASH, "/", nil)
				break;
			}
			for ; scanner.Peek() != '\n' && scanner.Peek() != 0 ; {
				scanner.Advance()
			}
			break;
		case '"':
			string_literal := ""
			for {
				if scanner.AtEnd() || scanner.Peek() == '\n' {
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
					found_error = true
					break
				}
				char := scanner.Advance()
				if char == '"'{
					scanner.AddToken(STRING, "\"" + string_literal + "\"", string_literal)
					break
				}
				string_literal += string(char)
			}
			
		case '\t':
		case ' ':
			break;
		default:
			if isDigit(char) {
				scanner.ScanNumber()
			} else if isAlpha(char) {
				for {
					peek := scanner.Peek()
					if !isAlphaNumeric(peek) { break }
					scanner.Advance()
				}
				lexeme := scanner.contents[scanner.start : scanner.current]
				if token_type, ok := reserved_words[lexeme]; ok {
					scanner.AddToken(token_type, lexeme, nil)
				} else {
					scanner.AddToken(IDENTIFIER, lexeme, nil)
				}
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, char)
				found_error = true
			}
		}
	}
	scanner.AddToken(EOF, "", nil)
	if found_error {
		return errors.New("Error scanning file contents")
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
	return ret;
}