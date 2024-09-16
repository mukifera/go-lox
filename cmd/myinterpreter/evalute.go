package main

import (
	"math/big"
	"fmt"
)

type Evaluator struct {
	expressions []Expression
	current int
	values []interface{}
}

func NewEvaluator(expressions []Expression) *Evaluator {
	var evaluator Evaluator
	evaluator.expressions = expressions
	evaluator.current = 0
	return &evaluator
}

func (evaluator *Evaluator) Evaluate() {
	evaluator.values = nil
	for _, expression := range evaluator.expressions {
		evaluator.values = append(evaluator.values, evaluator.evaluateExpression(expression))
	}
}

func (evaluator *Evaluator) evaluateExpression(expression Expression) interface{} {
	switch expression.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return expression.literal
	case ExpressionTypeEnum.GROUPING:
		return evaluator.evaluateExpression(expression.children[0])
	case ExpressionTypeEnum.UNARY:
		return evaluator.evaluateUnaryExpression(expression)
	case ExpressionTypeEnum.BINARY:
		return evaluator.evaluateBinaryExpression(expression)
	}
	return nil
}

func (evaluator *Evaluator) evaluateUnaryExpression(expression Expression) interface{} {
	value := evaluator.evaluateExpression(expression.children[0])
	switch expression.operator {
	case OperatorEnum.BANG:
		return value == false || value == nil
	case OperatorEnum.MINUS:
		number, ok := value.(big.Float)
		if !ok {
			break
		}
		return *number.Neg(&number)
	}
	return nil
}

func (evaluator *Evaluator) evaluateBinaryExpression(expression Expression) interface{} {
	left_value := evaluator.evaluateExpression(expression.children[0])
	right_value := evaluator.evaluateExpression(expression.children[1])
	switch expression.operator {
	case OperatorEnum.STAR:
		left, right, ok := evaluator.assertNumberOperation(left_value, right_value)
		if !ok {
			break
		}
		return *left.Mul(&left, &right)
	case OperatorEnum.SLASH:
		left, right, ok := evaluator.assertNumberOperation(left_value, right_value)
		if !ok {
			break
		}
		return *left.Quo(&left, &right)
	case OperatorEnum.MINUS:
		left, right, ok := evaluator.assertNumberOperation(left_value, right_value)
		if !ok {
			break
		}
		return *left.Sub(&left, &right)
	case OperatorEnum.PLUS:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return *left.Add(&left, &right)
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left + right
		}
		break
	case OperatorEnum.LESS:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) == -1
		}
		break
	case OperatorEnum.LESS_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) < 1
		}
		break
	case OperatorEnum.GREATER:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) == 1
		}
		break
	case OperatorEnum.GREATER_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) > -1
		}
		break
	case OperatorEnum.EQUAL_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) == 0
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left == right
		}
		return false
	case OperatorEnum.BANG_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) != 0
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left != right
		}
		return true
	}
	return nil
}

func (evaluator *Evaluator) assertNumberOperation(left_value interface{}, right_value interface{}) (big.Float, big.Float, bool) {
	left, left_ok := left_value.(big.Float)
	right, right_ok := right_value.(big.Float)
	return left, right, left_ok && right_ok
}

func (evaluator *Evaluator) assertStringOperation(left_value interface{}, right_value interface{}) (string, string, bool) {
	left, left_ok := left_value.(string)
	right, right_ok := right_value.(string)
	return left, right, left_ok && right_ok
}

func (evaluator *Evaluator) StringifyValues() []string {
	var ret []string
	for _, value := range evaluator.values {
		str := ""
		switch literal := value.(type)	{
		case nil:		 		str = fmt.Sprintf("nil")
		case bool:	 		str = fmt.Sprintf("%v", literal)
		case string: 		str = fmt.Sprintf("%s", literal)
		case big.Float: str = fmt.Sprintf("%s", literal.String())
		}
		ret = append(ret, str)
	}
	return ret
}