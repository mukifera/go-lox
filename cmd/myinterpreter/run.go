package main

import (
	"fmt"
	"io"
	"os"
)

func RunExpressions(exprs []Expression) []error {
	return FRunExpressions(os.Stdout, exprs)
}

func RunExpression(expr Expression, scope map[string]interface{}) error {
	return FRunExpression(os.Stdout, expr, scope)
}

func FRunExpressions(writer io.Writer, exprs []Expression) []error {
	errs := make([]error, len(exprs))
	scope := make(map[string]interface{})
	for index, expr := range exprs {
		err := FRunExpression(writer, expr, scope)
		errs[index] = err
		if err != nil {
			break
		}
	}
	return errs
}

func FRunExpression(writer io.Writer, expr Expression, scope map[string]interface{}) error {
	switch expr.operator {
	case OperatorEnum.PRINT:
		value, err := EvaluateExpression(expr.children[0], scope)
		if err != nil {
			return err
		}
		str := StringifyEvaluationValue(value)
		fmt.Fprintln(writer, str)
	case OperatorEnum.VAR:
		expr_identifier := expr.children[0].children[0]
		expr_value := expr.children[0].children[1]
		value, err := EvaluateExpression(expr_value, scope)
		if err != nil {
			return err
		}
		scope[expr_identifier.literal.(string)] = value

	default:
		_, err := EvaluateExpression(expr, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
