package main

import (
	"errors"
	"fmt"
	"os"
)

type Parser struct {
	current int
	has_error bool

	tokens []Token
	expressions []Expression
}

func NewParser(tokens []Token) *Parser {
	var parser Parser
	parser.current = 0
	parser.has_error = false
	parser.tokens = tokens
	return &parser
}

func (parser *Parser) ParseExpressions() []Expression {
	var expressions []Expression
	for {
		if parser.AtEnd() { break; }
		var expr Expression
		token := parser.Peek()
		switch token.token_type {
		case EOF, RIGHT_PAREN: return expressions
		case LEFT_PAREN:
			parser.Advance()
			expr.expression_type = ExpressionTypeEnum.GROUPING
			group := parser.ParseExpressions()
			if parser.Peek().token_type != RIGHT_PAREN {
				fmt.Fprintf(os.Stderr, "Error: Unmatched parentheses.\n")
				parser.has_error = true
				break
			}
			parser.Advance()
			expr.children = group
			break
		case TRUE:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = true
			parser.Advance()
			break
		case FALSE:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = false
			parser.Advance()
			break
		case NIL:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = nil
			parser.Advance()
			break
		case NUMBER, STRING:
			expr.expression_type = ExpressionTypeEnum.LITERAL
			expr.literal = token.literal
			parser.Advance()
			break
		}
		expressions = append(expressions, expr)
	}
	return expressions
}

func (parser *Parser) Parse() error {
	parser.expressions = parser.ParseExpressions()
	if parser.Peek().token_type != EOF {
		fmt.Fprintf(os.Stderr, "Error: Unmatched parentheses.\n")
		parser.has_error = true
	}
	if parser.has_error {
		return errors.New("Error parsing tokens")
	}
	return nil
}

func (parser *Parser) Advance() Token {
	token := parser.tokens[parser.current]
	parser.current += 1
	return token
}

func (parser *Parser) Peek() Token {
	if parser.AtEnd() {
		var token Token
		token.token_type = EOF
		return token
	}
	return parser.tokens[parser.current]
}

func (parser *Parser) AtEnd() bool {
	return parser.current >= len(parser.tokens) || parser.tokens[parser.current].token_type == EOF
}

func (parser *Parser) StringifyExpressions() string {
	str := ""
	for _, expr := range parser.expressions {
		str += expr.String()
	}
	return str
}