package main

import (
	"math/big"
	"fmt"
	"errors"
)

type numberBinaryOperation func(big.Float, big.Float) interface{}

type Evaluator struct {
	expressions []Expression
	current int
	values []interface{}
	errors []error
}

func NewEvaluator(expressions []Expression) *Evaluator {
	var evaluator Evaluator
	evaluator.expressions = expressions
	evaluator.current = 0
	return &evaluator
}

func (evaluator *Evaluator) Evaluate() {
	evaluator.values = nil
	evaluator.errors = nil
	for _, expression := range evaluator.expressions {
		value, err := evaluator.evaluateExpression(expression)
		evaluator.values = append(evaluator.values, value)
		evaluator.errors = append(evaluator.errors, err)
	}
}

func (evaluator *Evaluator) evaluateExpression(expression Expression) (interface{}, error) {
	switch expression.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return expression.literal, nil
	case ExpressionTypeEnum.GROUPING:
		return evaluator.evaluateExpression(expression.children[0])
	case ExpressionTypeEnum.UNARY:
		return evaluator.evaluateUnaryExpression(expression)
	case ExpressionTypeEnum.BINARY:
		return evaluator.evaluateBinaryExpression(expression)
	}
	return nil, nil
}

func (evaluator *Evaluator) evaluateUnaryExpression(expression Expression) (interface{}, error) {
	value, err := evaluator.evaluateExpression(expression.children[0])
	if err != nil {
		return nil, err
	}
	switch expression.operator {
	case OperatorEnum.BANG:
		return (value == false || value == nil), nil
	case OperatorEnum.MINUS:
		number, ok := value.(big.Float)
		if !ok {
			return nil, errors.New("Operand must be a number.")
		}
		return *number.Neg(&number), nil
	}
	return nil, errors.New("Unknown unary operator.")
}

func (evaluator *Evaluator) evaluateBinaryExpression(expression Expression) (interface{}, error) {
	left_value, err := evaluator.evaluateExpression(expression.children[0])
	if err != nil {
		return nil, err
	}
	right_value, err := evaluator.evaluateExpression(expression.children[1])
	if err != nil {
		return nil, err
	}

	num_or_str_operation_error := errors.New("Operands must be two numbers or two strings.")

	exec_number_operation := func(operation numberBinaryOperation) (interface{}, error) {
		left, right, ok := evaluator.assertNumberOperation(left_value, right_value)
		if !ok {
			return nil, errors.New("Operands must be numbers.")
		}
		return operation(left, right), nil
	}

	switch expression.operator {
	case OperatorEnum.STAR:
		return exec_number_operation(func(l big.Float, r big.Float) interface{} {
			return *l.Mul(&l, &r)
		})
	case OperatorEnum.SLASH:
		return exec_number_operation(func(l big.Float, r big.Float) interface{} {
			return *l.Quo(&l, &r)
		})
	case OperatorEnum.MINUS:
		return exec_number_operation(func(l big.Float, r big.Float) interface{} {
			return *l.Sub(&l, &r)
		})
	case OperatorEnum.PLUS:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return *left.Add(&left, &right), nil
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left + right, nil
		}
		return nil, num_or_str_operation_error
	case OperatorEnum.LESS:
		return exec_number_operation(func(left big.Float, right big.Float) interface{} {
			return left.Cmp(&right) == -1
		})
	case OperatorEnum.LESS_EQUAL:
		return exec_number_operation(func(left big.Float, right big.Float) interface{} {
			return left.Cmp(&right) < 1
		})
	case OperatorEnum.GREATER:
		return exec_number_operation(func(left big.Float, right big.Float) interface{} {
			return left.Cmp(&right) == 1
		})
	case OperatorEnum.GREATER_EQUAL:
		return exec_number_operation(func(left big.Float, right big.Float) interface{} {
			return left.Cmp(&right) > -1
		})
	case OperatorEnum.EQUAL_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) == 0, nil
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left == right, nil
		}
		return false, nil
	case OperatorEnum.BANG_EQUAL:
		if left, right, ok := evaluator.assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) != 0, nil
		}
		if left, right, ok := evaluator.assertStringOperation(left_value, right_value); ok {
			return left != right, nil
		}
		return true, nil
	}
	return nil, errors.New("Unkown binary operator.")
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