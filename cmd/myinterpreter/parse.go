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

func (parser *Parser) ParseOneExpression() Expression {
	var expr Expression
	token := parser.Peek()
	switch token.token_type {
	case EOF, RIGHT_PAREN: return expr
	case LEFT_PAREN:
		parser.Advance()
		expr.expression_type = ExpressionTypeEnum.GROUPING
		group := parser.ParseExpressions()
		if parser.Advance().token_type != RIGHT_PAREN {
			fmt.Fprintf(os.Stderr, "Error: Unmatched parentheses.\n")
			parser.has_error = true
			break
		}
		if len(group) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Parentheses contain no expression.\n")
			parser.has_error = true
			break
		}
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
	case BANG, MINUS:
		expr.expression_type = ExpressionTypeEnum.UNARY
		expr.operator = OperatorEnum.MINUS
		if token.token_type == BANG {
			expr.operator = OperatorEnum.BANG
		}
		expr.literal = nil
		parser.Advance()
		child := parser.ParseOneExpression()
		if child.expression_type == ExpressionTypeEnum.UNDEFINED {
			fmt.Fprintf(os.Stderr, "Error: Expected an expression after unary operator %s.\n", parser.Peek().StringLiteral())
			parser.has_error = true
			break
		}
		expr.children = append(expr.children, child)
		break
	case STAR, SLASH:
		expr.expression_type = ExpressionTypeEnum.BINARY
		expr.operator = OperatorEnum.STAR
		if token.token_type == SLASH {
			expr.operator = OperatorEnum.SLASH
		}
		expr.literal = nil
		parser.Advance()
		break
	}
	return expr
}

func (parser *Parser) ParseExpressions() []Expression {
	var expressions []Expression
	for {
		if parser.AtEnd() { break; }
		expr := parser.ParseOneExpression()
		if expr.expression_type == ExpressionTypeEnum.UNDEFINED {
			break
		}
		expressions = append(expressions, expr)
	}
	for i := 0; i < len(expressions); {
		expr := expressions[i]
		if expr.expression_type == ExpressionTypeEnum.BINARY {
			if i == 0 {
				fmt.Fprintf(os.Stderr, "Error: Expected an expression before binary operator %s.\n", parser.Peek().StringLiteral())
				parser.has_error = true
				break
			}
			if i + 1 >= len(expressions) {
				fmt.Fprintf(os.Stderr, "Error: Expected an expression after binary operator %s.\n", parser.Peek().StringLiteral())
				parser.has_error = true
				break
			}
			expr.children = append(expr.children, expressions[i-1], expressions[i+1])
			suffix := expressions[i+2:]
			expressions = append(expressions[:i-1], expr)
			expressions = append(expressions, suffix...)
			continue
		}
		i += 1
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