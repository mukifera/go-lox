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

// func (parser *Parser) ParseOneExpression() Expression {
// 	var expr Expression
// 	token := parser.Peek()
// 	switch token.token_type {
// 	case EOF, RIGHT_PAREN: return expr
// 	case LEFT_PAREN:
// 		parser.Advance()
// 		expr.expression_type = ExpressionTypeEnum.GROUPING
// 		group := parser.ParseExpressions()
// 		if parser.Advance().token_type != RIGHT_PAREN {
// 			fmt.Fprintf(os.Stderr, "Error: Unmatched parentheses.\n")
// 			parser.has_error = true
// 			break
// 		}
// 		if len(group) == 0 {
// 			fmt.Fprintf(os.Stderr, "Error: Parentheses contain no expression.\n")
// 			parser.has_error = true
// 			break
// 		}
// 		expr.children = group
// 		break
// 	case TRUE:
// 		expr.expression_type = ExpressionTypeEnum.LITERAL;
// 		expr.literal = true
// 		parser.Advance()
// 		break
// 	case FALSE:
// 		expr.expression_type = ExpressionTypeEnum.LITERAL;
// 		expr.literal = false
// 		parser.Advance()
// 		break
// 	case NIL:
// 		expr.expression_type = ExpressionTypeEnum.LITERAL;
// 		expr.literal = nil
// 		parser.Advance()
// 		break
// 	case NUMBER, STRING:
// 		expr.expression_type = ExpressionTypeEnum.LITERAL
// 		expr.literal = token.literal
// 		parser.Advance()
// 		break
// 	case BANG:
// 		expr.expression_type = ExpressionTypeEnum.UNARY
// 		expr.operator = OperatorEnum.BANG
// 		expr.literal = nil
// 		parser.Advance()
// 		break
// 	case STAR, SLASH, PLUS, MINUS, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL, EQUAL_EQUAL, BANG_EQUAL:
// 		expr.expression_type = ExpressionTypeEnum.BINARY
// 		expr.operator = OperatorEnum.STAR
// 		if token.token_type == SLASH {
// 			expr.operator = OperatorEnum.SLASH
// 		} else if token.token_type == PLUS {
// 			expr.operator = OperatorEnum.PLUS
// 		} else if token.token_type == MINUS {
// 			expr.operator = OperatorEnum.MINUS
// 		} else if token.token_type == LESS {
// 			expr.operator = OperatorEnum.LESS
// 		} else if token.token_type == LESS_EQUAL {
// 			expr.operator = OperatorEnum.LESS_EQUAL
// 		} else if token.token_type == GREATER {
// 			expr.operator = OperatorEnum.GREATER
// 		} else if token.token_type == GREATER_EQUAL {
// 			expr.operator = OperatorEnum.GREATER_EQUAL
// 		} else if token.token_type == EQUAL_EQUAL {
// 			expr.operator = OperatorEnum.EQUAL_EQUAL
// 		} else if token.token_type == BANG_EQUAL {
// 			expr.operator = OperatorEnum.BANG_EQUAL
// 		}
// 		expr.literal = nil
// 		parser.Advance()
// 		break
// 	}
// 	return expr
// }

// func (parser *Parser) UpdateBinaryExpressions(expressions []Expression, operators []Operator) []Expression {
// 	var ret []Expression
// 	for i := 0; i < len(expressions); i++ {
// 		expr := expressions[i]
// 		if expr.expression_type != ExpressionTypeEnum.BINARY {
// 			ret = append(ret, expr)
// 			continue
// 		}
// 		found := false
// 		for _, operator := range operators {
// 			if expr.operator == operator {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			ret = append(ret, expr)
// 			continue
// 		}
// 		if len(ret) == 0 {
// 			fmt.Fprintf(os.Stderr, "Error: Expected an expression before binary operator %s.\n", parser.Peek().StringLiteral())
// 			parser.has_error = true
// 			break
// 		}
// 		if i + 1 >= len(expressions) {
// 			fmt.Fprintf(os.Stderr, "Error: Expected an expression after binary operator %s.\n", parser.Peek().StringLiteral())
// 			parser.has_error = true
// 			break
// 		}
// 		expr.children = append(expr.children, ret[len(ret)-1], expressions[i+1])
// 		ret[len(ret) - 1] = expr
// 		i += 1
// 	}
// 	return ret
// }

// func (parser *Parser) ParseExpressions() []Expression {
// 	var expressions []Expression
// 	for {
// 		if parser.AtEnd() { break; }
// 		expr := parser.ParseOneExpression()
// 		if expr.expression_type == ExpressionTypeEnum.UNDEFINED {
// 			break
// 		}
// 		expressions = append(expressions, expr)
// 	}

// 	if len(expressions) > 0 &&
// 		expressions[0].operator == OperatorEnum.MINUS {
// 			expressions[0].expression_type = ExpressionTypeEnum.UNARY
// 	}
// 	for i := 0; i + 1 < len(expressions); i++ {
// 		current_type := expressions[i].expression_type
// 		next_type := expressions[i+1].expression_type
		
// 		if (current_type == ExpressionTypeEnum.BINARY || current_type == ExpressionTypeEnum.UNARY) &&
// 			(next_type == ExpressionTypeEnum.BINARY) &&
// 			(expressions[i+1].operator == OperatorEnum.MINUS) {
// 			expressions[i+1].expression_type = ExpressionTypeEnum.UNARY
// 		}
// 	}

// 	for i := len(expressions) - 1; i >= 0; i-- {
// 		if expressions[i].expression_type == ExpressionTypeEnum.UNARY {
// 			if i + 1 >= len(expressions) {
// 				fmt.Fprintf(os.Stderr, "Error: Expected an expression after unary operator %s.\n", parser.Peek().StringLiteral())
// 				parser.has_error = true
// 				break
// 			}
// 			expressions[i].children = append(expressions[i].children, expressions[i+1])
// 			expressions = append(expressions[:i+1], expressions[i+2:]...)
// 		}
// 	}

// 	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.STAR, OperatorEnum.SLASH})
// 	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.PLUS, OperatorEnum.MINUS})
// 	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.LESS, OperatorEnum.LESS_EQUAL, OperatorEnum.GREATER, OperatorEnum.GREATER_EQUAL})
// 	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.EQUAL_EQUAL, OperatorEnum.BANG_EQUAL})
	
// 	return expressions
// }

func (parser *Parser) Parse() error {
	expr := parser.parseExpression()
	parser.expressions = []Expression{expr}
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
	if parser.Matches(NUMBER, STRING) {
		return NewLiteralExpression(parser.Previous().literal)
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