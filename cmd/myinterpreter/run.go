package main

import (
	"fmt"
	"io"
	"os"
)

func RunExpressions(exprs []Expression) []error {
	return FRunExpressions(os.Stdout, exprs)
}

func RunExpression(expr Expression) error {
	return FRunExpression(os.Stdout, expr)
}

func FRunExpressions(writer io.Writer, exprs []Expression) []error {
	errs := make([]error, len(exprs))
	for index, expr := range exprs {
		err := FRunExpression(writer, expr)
		errs[index] = err
		if err != nil {
			break
		}
	}
	return errs
}

func FRunExpression(writer io.Writer, expr Expression) error {
	switch expr.operator {
	case OperatorEnum.PRINT:
		value, err := EvaluateExpression(expr.children[0])
		if err != nil {
			return err
		}
		str := StringifyEvaluationValue(value)
		fmt.Fprintln(writer, str)
	default:
		_, err := EvaluateExpression(expr)
		if err != nil {
			return err
		}
	}
	return nil
}
