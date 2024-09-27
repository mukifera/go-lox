package main

import (
	"errors"
	"fmt"
	"math/big"
)

type numberBinaryOperation func(big.Float, big.Float) interface{}

func EvaluateExpressions(exprs []Expression) ([]interface{}, []error) {
	var scope map[string]interface{}
	values := make([]interface{}, len(exprs))
	errs := make([]error, len(exprs))
	for index, expr := range exprs {
		value, err := EvaluateExpression(expr, scope)
		values[index] = value
		errs[index] = err
	}
	return values, errs
}

func EvaluateExpression(expr Expression, scope map[string]interface{}) (interface{}, error) {
	switch expr.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return expr.literal, nil
	case ExpressionTypeEnum.GROUPING:
		return EvaluateExpression(expr.children[0], scope)
	case ExpressionTypeEnum.UNARY:
		return evaluateUnaryExpression(expr, scope)
	case ExpressionTypeEnum.BINARY:
		return evaluateBinaryExpression(expr, scope)
	case ExpressionTypeEnum.IDENTIFIER:
		variable := expr.literal.(string)
		value, ok := scope[variable]
		if !ok {
			return nil, fmt.Errorf("undefined variable '%s'", variable)
		}
		return value, nil
	}
	return nil, nil
}

func evaluateUnaryExpression(expr Expression, scope map[string]interface{}) (interface{}, error) {
	value, err := EvaluateExpression(expr.children[0], scope)
	if err != nil {
		return nil, err
	}
	switch expr.operator {
	case OperatorEnum.BANG:
		return (value == false || value == nil), nil
	case OperatorEnum.MINUS:
		number, ok := value.(big.Float)
		if !ok {
			return nil, errors.New("operand must be a number")
		}
		return *number.Neg(&number), nil
	}
	return nil, errors.New("unknown unary operator")
}

func evaluateBinaryExpression(expr Expression, scope map[string]interface{}) (interface{}, error) {
	left_value, err := EvaluateExpression(expr.children[0], scope)
	if err != nil {
		return nil, err
	}
	right_value, err := EvaluateExpression(expr.children[1], scope)
	if err != nil {
		return nil, err
	}

	num_or_str_operation_error := errors.New("operands must be two numbers or two strings")

	exec_number_operation := func(operation numberBinaryOperation) (interface{}, error) {
		left, right, ok := assertNumberOperation(left_value, right_value)
		if !ok {
			return nil, errors.New("operands must be numbers")
		}
		return operation(left, right), nil
	}

	switch expr.operator {
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
		if left, right, ok := assertNumberOperation(left_value, right_value); ok {
			return *left.Add(&left, &right), nil
		}
		if left, right, ok := assertStringOperation(left_value, right_value); ok {
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
		if left, right, ok := assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) == 0, nil
		}
		if left, right, ok := assertStringOperation(left_value, right_value); ok {
			return left == right, nil
		}
		if left, right, ok := assertBoolOperation(left_value, right_value); ok {
			return left == right, nil
		}
		return false, nil
	case OperatorEnum.BANG_EQUAL:
		if left, right, ok := assertNumberOperation(left_value, right_value); ok {
			return left.Cmp(&right) != 0, nil
		}
		if left, right, ok := assertStringOperation(left_value, right_value); ok {
			return left != right, nil
		}
		if left, right, ok := assertBoolOperation(left_value, right_value); ok {
			return left != right, nil
		}
		return true, nil
	}
	return nil, errors.New("unkown binary operator")
}

func assertNumberOperation(left_value interface{}, right_value interface{}) (big.Float, big.Float, bool) {
	left, left_ok := left_value.(big.Float)
	right, right_ok := right_value.(big.Float)
	return left, right, left_ok && right_ok
}

func assertStringOperation(left_value interface{}, right_value interface{}) (string, string, bool) {
	left, left_ok := left_value.(string)
	right, right_ok := right_value.(string)
	return left, right, left_ok && right_ok
}

func assertBoolOperation(left_value interface{}, right_value interface{}) (bool, bool, bool) {
	left, left_ok := left_value.(bool)
	right, right_ok := right_value.(bool)
	return left, right, left_ok && right_ok
}

func StringifyEvaluationValues(values []interface{}) []string {
	strs := make([]string, len(values))
	for index, value := range values {
		str := StringifyEvaluationValue(value)
		strs[index] = str
	}
	return strs
}

func StringifyEvaluationValue(value interface{}) string {
	switch literal := value.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%v", literal)
	case string:
		return literal
	case big.Float:
		return literal.String()
	}
	return ""
}
