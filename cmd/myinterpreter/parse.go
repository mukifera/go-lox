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

func tokenTypeToOperator(token_type TokenType) Operator {
	switch token_type {
		case MINUS: 				return OperatorEnum.MINUS
		case PLUS: 					return OperatorEnum.PLUS
		case STAR: 					return OperatorEnum.STAR
		case EQUAL: 				return OperatorEnum.EQUAL
		case EQUAL_EQUAL: 	return OperatorEnum.EQUAL_EQUAL
		case BANG: 					return OperatorEnum.BANG
		case BANG_EQUAL: 		return OperatorEnum.BANG_EQUAL
		case LESS: 					return OperatorEnum.LESS
		case LESS_EQUAL: 		return OperatorEnum.LESS_EQUAL
		case GREATER: 			return OperatorEnum.GREATER
		case GREATER_EQUAL: return OperatorEnum.GREATER_EQUAL
		case SLASH: 				return OperatorEnum.SLASH
		}
	return OperatorEnum.UNDEFINED
}

func (parser *Parser) Parse() error {
	for !parser.AtEnd() {
		expr := parser.parseExpression()
		if !parser.AtEnd() && !parser.Matches(SEMICOLON) {
			fmt.Fprintf(os.Stderr, "Error: Expected semicolon.\n")
			parser.has_error = true
		}
		parser.expressions = append(parser.expressions, expr)
	}
	if parser.has_error {
		return errors.New("Error parsing tokens")
	}
	return nil
}

func (parser *Parser) parseExpression() Expression {
	return parser.parseEquality()
}

func (parser *Parser) parseEquality() Expression {
	expr := parser.parseComparison()

	for parser.Matches(EQUAL_EQUAL, BANG_EQUAL) {
		token_type := parser.Previous().token_type
		operator := tokenTypeToOperator(token_type)
		right := parser.parseComparison()
		top := NewBinaryExpression(expr, right, operator)
		expr = top
	}
	return expr
}

func (parser *Parser) parseComparison() Expression {
	var expr Expression = parser.parseAddSub()

	for parser.Matches(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression = parser.parseAddSub()
		var top Expression = NewBinaryExpression(expr, right, operator)
		expr = top
	}
	return expr
}

func (parser *Parser) parseAddSub() Expression {
	var expr Expression = parser.parseMultDiv()

	for parser.Matches(PLUS, MINUS) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression = parser.parseMultDiv()
		var top Expression = NewBinaryExpression(expr, right, operator)
		expr = top
	}
	return expr
}

func (parser *Parser) parseMultDiv() Expression {
	var expr Expression = parser.parseUnary()

	for parser.Matches(STAR, SLASH) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression = parser.parseUnary()
		var top Expression = NewBinaryExpression(expr, right, operator)
		expr = top
	}
	return expr
}

func (parser *Parser) parseUnary() Expression {
	if parser.Matches(BANG, MINUS) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var expr Expression = parser.parseUnary()
		return NewUnaryExpression(expr, operator)
	}
	return parser.parsePrimary()
}

func (parser *Parser) parsePrimary() Expression {
	if parser.Matches(FALSE) {
		return NewLiteralExpression(false)
	}
	if parser.Matches(TRUE) {
		return NewLiteralExpression(true)
	}
	if parser.Matches(NIL) {
		return NewLiteralExpression(nil)
	}
	if parser.Matches(NUMBER, STRING, IDENTIFIER) {
		return NewLiteralExpression(parser.Previous().literal)
	}
	if parser.Matches(PRINT) {
		expr := parser.parseExpression()
		return NewBuiltinExpression(expr, OperatorEnum.PRINT)
	}
	if parser.Matches(LEFT_PAREN) {
		expr := parser.parseExpression()
		if !parser.Matches(RIGHT_PAREN) {
			fmt.Fprintf(os.Stderr, "Error: Unmatched parentheses.\n")
			parser.has_error = true
		}
		return NewGroupingExpression(expr)
	}
	fmt.Fprintf(os.Stderr, "Error: Unknown Token.\n")
	parser.has_error = true
	return NewUndefinedExpression()
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

func (parser *Parser) Matches(token_types ...TokenType) bool {
	for _, token_type := range token_types {
		if token_type == parser.Peek().token_type {
			parser.Advance()
			return true
		}
	}
	return false
}

func (parser *Parser) Previous() Token {
	return parser.tokens[parser.current - 1]
}

func (parser *Parser) StringifyExpressions() string {
	str := ""
	for _, expr := range parser.expressions {
		str += expr.String()
	}
	return str
}