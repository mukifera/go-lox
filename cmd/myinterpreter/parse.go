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
	case BANG:
		expr.expression_type = ExpressionTypeEnum.UNARY
		expr.operator = OperatorEnum.BANG
		expr.literal = nil
		parser.Advance()
		break
	case STAR, SLASH, PLUS, MINUS, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL, EQUAL_EQUAL, BANG_EQUAL:
		expr.expression_type = ExpressionTypeEnum.BINARY
		expr.operator = OperatorEnum.STAR
		if token.token_type == SLASH {
			expr.operator = OperatorEnum.SLASH
		} else if token.token_type == PLUS {
			expr.operator = OperatorEnum.PLUS
		} else if token.token_type == MINUS {
			expr.operator = OperatorEnum.MINUS
		} else if token.token_type == LESS {
			expr.operator = OperatorEnum.LESS
		} else if token.token_type == LESS_EQUAL {
			expr.operator = OperatorEnum.LESS_EQUAL
		} else if token.token_type == GREATER {
			expr.operator = OperatorEnum.GREATER
		} else if token.token_type == GREATER_EQUAL {
			expr.operator = OperatorEnum.GREATER_EQUAL
		} else if token.token_type == EQUAL_EQUAL {
			expr.operator = OperatorEnum.EQUAL_EQUAL
		} else if token.token_type == BANG_EQUAL {
			expr.operator = OperatorEnum.BANG_EQUAL
		}
		expr.literal = nil
		parser.Advance()
		break
	}
	return expr
}

func (parser *Parser) UpdateBinaryExpressions(expressions []Expression, operators []Operator) []Expression {
	var ret []Expression
	for i := 0; i < len(expressions); i++ {
		expr := expressions[i]
		if expr.expression_type != ExpressionTypeEnum.BINARY {
			ret = append(ret, expr)
			continue
		}
		found := false
		for _, operator := range operators {
			if expr.operator == operator {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, expr)
			continue
		}
		if len(ret) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Expected an expression before binary operator %s.\n", parser.Peek().StringLiteral())
			parser.has_error = true
			break
		}
		if i + 1 >= len(expressions) {
			fmt.Fprintf(os.Stderr, "Error: Expected an expression after binary operator %s.\n", parser.Peek().StringLiteral())
			parser.has_error = true
			break
		}
		expr.children = append(expr.children, ret[len(ret)-1], expressions[i+1])
		ret[len(ret) - 1] = expr
		i += 1
	}
	return ret
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

	if len(expressions) > 0 &&
		expressions[0].operator == OperatorEnum.MINUS {
			expressions[0].expression_type = ExpressionTypeEnum.UNARY
	}
	for i := 0; i + 1 < len(expressions); i++ {
		current_type := expressions[i].expression_type
		next_type := expressions[i+1].expression_type
		
		if (current_type == ExpressionTypeEnum.BINARY || current_type == ExpressionTypeEnum.UNARY) &&
			(next_type == ExpressionTypeEnum.BINARY) &&
			(expressions[i+1].operator == OperatorEnum.MINUS) {
			expressions[i+1].expression_type = ExpressionTypeEnum.UNARY
		}
	}

	for i := len(expressions) - 1; i >= 0; i-- {
		if expressions[i].expression_type == ExpressionTypeEnum.UNARY {
			if i + 1 >= len(expressions) {
				fmt.Fprintf(os.Stderr, "Error: Expected an expression after unary operator %s.\n", parser.Peek().StringLiteral())
				parser.has_error = true
				break
			}
			expressions[i].children = append(expressions[i].children, expressions[i+1])
			expressions = append(expressions[:i+1], expressions[i+2:]...)
		}
	}

	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.STAR, OperatorEnum.SLASH})
	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.PLUS, OperatorEnum.MINUS})
	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.LESS, OperatorEnum.LESS_EQUAL, OperatorEnum.GREATER, OperatorEnum.GREATER_EQUAL})
	expressions = parser.UpdateBinaryExpressions(expressions, []Operator{OperatorEnum.EQUAL_EQUAL, OperatorEnum.BANG_EQUAL})
	
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