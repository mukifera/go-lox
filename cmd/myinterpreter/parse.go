package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	current int

	tokens      []Token
	expressions []Expression
}

func NewParser(tokens []Token) *Parser {
	var parser Parser
	parser.current = 0
	parser.tokens = tokens
	return &parser
}

func tokenTypeToOperator(token_type TokenType) Operator {
	switch token_type {
	case MINUS:
		return OperatorEnum.MINUS
	case PLUS:
		return OperatorEnum.PLUS
	case STAR:
		return OperatorEnum.STAR
	case EQUAL:
		return OperatorEnum.EQUAL
	case EQUAL_EQUAL:
		return OperatorEnum.EQUAL_EQUAL
	case BANG:
		return OperatorEnum.BANG
	case BANG_EQUAL:
		return OperatorEnum.BANG_EQUAL
	case LESS:
		return OperatorEnum.LESS
	case LESS_EQUAL:
		return OperatorEnum.LESS_EQUAL
	case GREATER:
		return OperatorEnum.GREATER
	case GREATER_EQUAL:
		return OperatorEnum.GREATER_EQUAL
	case SLASH:
		return OperatorEnum.SLASH
	case OR:
		return OperatorEnum.OR
	case AND:
		return OperatorEnum.AND
	}
	return OperatorEnum.UNDEFINED
}

func (parser *Parser) Parse() error {
	var err error
	parser.expressions, err = parser.parseStatements()
	return err
}

func (parser *Parser) parseStatements() ([]Expression, error) {
	exprs := make([]Expression, 0)
	var err error = nil
	for !parser.AtEnd() {
		expr, sub_err := parser.parseStatement()
		err = errors.Join(err, sub_err)
		exprs = append(exprs, expr)
		if parser.Peek().token_type == RIGHT_BRACE {
			break
		}
	}
	return exprs, err
}

func (parser *Parser) parseStatement() (Expression, error) {
	var expr Expression
	var sub_err error
	var sub_exprs []Expression
	var err error = nil
	if parser.Matches(IF) {
		var children []Expression
		sub, err := parser.parseExpression()
		if sub.expression_type != ExpressionTypeEnum.GROUPING {
			err = errors.Join(err, newParsingError("Error: Invalid if condition"))
		}
		children = append(children, sub)

		sub, sub_err := parser.parseStatement()
		err = errors.Join(err, sub_err)
		children = append(children, sub)

		if parser.Matches(ELSE) {
			sub, sub_err = parser.parseStatement()
			err = errors.Join(err, sub_err)
			children = append(children, sub)
		}
		return NewBuiltinExpression(OperatorEnum.IF, children...), err
	}
	if parser.Matches(LEFT_BRACE) {
		sub_exprs, sub_err = parser.parseStatements()
		err = errors.Join(err, sub_err)
		if !parser.Matches(RIGHT_BRACE) {
			err = errors.Join(err, newParsingError("Error: Unmatched curly brace"))
		}
		expr = NewScopeExpression(sub_exprs...)
		return expr, err
	}
	expr, sub_err = parser.parseExpression()
	err = errors.Join(err, sub_err)
	if !parser.AtEnd() && !parser.Matches(SEMICOLON) {
		fmt.Println(expr.String())
		err = errors.Join(err, newParsingError("Error: Expected semicolon"))
	}
	return expr, err
}

func (parser *Parser) parseExpression() (Expression, error) {
	return parser.parseAssignment()
}

func (parser *Parser) parseAssignment() (Expression, error) {
	expr, err := parser.parseEquality()

	exprs := []Expression{expr}

	for parser.Matches(EQUAL) {
		right, sub_err := parser.parseEquality()
		exprs = append(exprs, right)
		err = errors.Join(err, sub_err)
	}

	for i := len(exprs) - 2; i >= 0; i-- {
		right := exprs[i+1]
		top := NewBinaryExpression(exprs[i], right, OperatorEnum.EQUAL)
		exprs[i] = top
	}

	return exprs[0], err
}

func (parser *Parser) parseEquality() (Expression, error) {
	expr, err := parser.parseLogical()

	for parser.Matches(EQUAL_EQUAL, BANG_EQUAL) {
		token_type := parser.Previous().token_type
		operator := tokenTypeToOperator(token_type)
		right, sub_err := parser.parseLogical()
		top := NewBinaryExpression(expr, right, operator)
		expr = top
		err = errors.Join(err, sub_err)
	}
	return expr, err
}

func (parser *Parser) parseLogical() (Expression, error) {
	expr, err := parser.parseComparison()

	for parser.Matches(OR, AND) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression
		var top Expression
		right, sub_err := parser.parseComparison()
		top = NewBinaryExpression(expr, right, operator)
		expr = top
		err = errors.Join(err, sub_err)
	}
	return expr, err
}

func (parser *Parser) parseComparison() (Expression, error) {
	expr, err := parser.parseAddSub()

	for parser.Matches(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression
		var top Expression
		right, sub_err := parser.parseAddSub()
		top = NewBinaryExpression(expr, right, operator)
		expr = top
		err = errors.Join(err, sub_err)
	}
	return expr, err
}

func (parser *Parser) parseAddSub() (Expression, error) {
	expr, err := parser.parseMultDiv()

	for parser.Matches(PLUS, MINUS) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression
		var top Expression
		right, sub_err := parser.parseMultDiv()
		top = NewBinaryExpression(expr, right, operator)
		expr = top
		err = errors.Join(err, sub_err)
	}
	return expr, err
}

func (parser *Parser) parseMultDiv() (Expression, error) {
	expr, err := parser.parseUnary()

	for parser.Matches(STAR, SLASH) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var right Expression
		var top Expression
		right, sub_err := parser.parseUnary()
		top = NewBinaryExpression(expr, right, operator)
		expr = top
		err = errors.Join(err, sub_err)
	}
	return expr, err
}

func (parser *Parser) parseUnary() (Expression, error) {
	if parser.Matches(BANG, MINUS) {
		var token_type TokenType = parser.Previous().token_type
		var operator Operator = tokenTypeToOperator(token_type)
		var expr Expression
		expr, err := parser.parseUnary()
		return NewUnaryExpression(expr, operator), err
	}
	return parser.parsePrimary()
}

func (parser *Parser) parsePrimary() (Expression, error) {
	if parser.Matches(FALSE) {
		return NewLiteralExpression(false), nil
	}
	if parser.Matches(TRUE) {
		return NewLiteralExpression(true), nil
	}
	if parser.Matches(NIL) {
		return NewLiteralExpression(nil), nil
	}
	if parser.Matches(NUMBER, STRING) {
		return NewLiteralExpression(parser.Previous().literal), nil
	}
	if parser.Matches(IDENTIFIER) {
		return NewIdentifierExpression(parser.Previous().lexeme), nil
	}
	if parser.Matches(PRINT) {
		expr, err := parser.parseExpression()
		return NewBuiltinExpression(OperatorEnum.PRINT, expr), err
	}
	if parser.Matches(VAR) {
		expr, err := parser.parseExpression()
		if expr.expression_type != ExpressionTypeEnum.IDENTIFIER &&
			(expr.expression_type != ExpressionTypeEnum.BINARY ||
				expr.operator != OperatorEnum.EQUAL ||
				expr.children[0].expression_type != ExpressionTypeEnum.IDENTIFIER) {
			err = errors.Join(err, newParsingError("Error: Invalid variable declaration"))
		}
		return NewBuiltinExpression(OperatorEnum.VAR, expr), err
	}
	if parser.Matches(LEFT_PAREN) {
		expr, err := parser.parseExpression()
		if !parser.Matches(RIGHT_PAREN) {
			err = errors.Join(err, newParsingError("Error: Unmatched parentheses"))
		}
		return NewGroupingExpression(expr), err
	}
	return NewUndefinedExpression(), newParsingError("Error: Unknown Token")
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
	return parser.tokens[parser.current-1]
}

func (parser *Parser) StringifyExpressions() string {
	str := ""
	for _, expr := range parser.expressions {
		str += expr.String() + "\n"
	}
	return str
}
