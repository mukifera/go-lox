package main

import "testing"

func getExpressions(fileContents string, t *testing.T) []Expression {
	scanner := NewScanner(fileContents)
	errs := scanner.Scan()
	if len(errs) != 0 {
		t.Errorf("Scanner: tokenizing error: %v", errs)
	}

	parser := NewParser(scanner.tokens)
	err := parser.Parse()
	if err != nil {
		t.Errorf("Parser: parsing error: %v", err)
	}

	return parser.expressions
}

func TestEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		expected     []string
	}{
		{"Literals: Boolean/true", "true", []string{"true"}},
		{"Literals: Boolean/false", "false", []string{"false"}},
		{"Literals: nil", "nil", []string{"nil"}},
		{"Literals: Number #1", "10.40", []string{"10.4"}},
		{"Literals: Number #2", "10", []string{"10"}},
		{"Parentheses #1", `("hello world!")`, []string{"hello world!"}},
		{"Parentheses #2", "(true)", []string{"true"}},
		{"Parentheses #3", "(10.40)", []string{"10.4"}},
		{"Parentheses #4", "((false))", []string{"false"}},
		{"Unary: Negation", "-73", []string{"-73"}},
		{"Unary: Not #1", "!true", []string{"false"}},
		{"Unary: Not #2", "!10.40", []string{"false"}},
		{"Unary: Not #3", "!((false))", []string{"true"}},
		{"Arithmetic #1", "42 / 5", []string{"8.4"}},
		{"Arithmetic #2", "18 * 3 / (3 * 6)", []string{"3"}},
		{"Arithmetic #3", "(10.40 * 2) / 2", []string{"10.4"}},
		{"Arithmetic #4", "70 - 65", []string{"5"}},
		{"Arithmetic #5", "69 - 93", []string{"-24"}},
		{"Arithmetic #6", "10.40 - 2", []string{"8.4"}},
		{"Arithmetic #6", "23 + 28 - (-(61 - 99))", []string{"13"}},
		{"String Concatenation #1", `"hello" + " world!"`, []string{"hello world!"}},
		{"String Concatenation #2", `"42" + "24"`, []string{"4224"}},
		{"String Concatenation #3", `"foo" + "bar"`, []string{"foobar"}},
		{"Relational Operators #1", "57 > -65", []string{"true"}},
		{"Relational Operators #2", "11 >= 11", []string{"true"}},
		{"Relational Operators #3", "(54 - 67) >= -(114 / 57 + 11)", []string{"true"}},
		{"Equality #1", `"hello" == "world"`, []string{"false"}},
		{"Equality #2", `"foo" != "bar"`, []string{"true"}},
		{"Equality #3", `"foo" == "foo"`, []string{"true"}},
		{"Equality #4", `61 == "61"`, []string{"false"}},
		{"Equality #5", "61 == 61", []string{"true"}},
		{"Equality #6", "61 == 10.5", []string{"false"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprs := getExpressions(tt.fileContents, t)
			values, _ := EvaluateExpressions(exprs)
			actual := StringifyEvaluationValues(values)

			if len(tt.expected) != len(actual) {
				t.Errorf("Evaluation result length mismatch: Expected %d outputs, got %d", len(tt.expected), len(actual))
			} else {
				for index, str := range actual {
					if tt.expected[index] != str {
						t.Errorf("Evaluation result mismatch on output %d\nExpected: %s\nGot: %s", index+1, tt.expected[index], str)
					}
				}
			}
		})
	}
}

func TestEvaluationRuntimeErrors(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		expected     []string
	}{
		{"Negation #1", `-"foo"`, []string{"operand must be a number"}},
		{"Negation #2", "-true", []string{"operand must be a number"}},
		{"Negation #3", `-("foo" + "bar")`, []string{"operand must be a number"}},
		{"Negation #4", "-false", []string{"operand must be a number"}},
		{"Multiplication #1", `"foo" * 42`, []string{"operands must be numbers"}},
		{"Multiplication #2", `("foo" * "bar")`, []string{"operands must be numbers"}},
		{"Division #1", "true / 2", []string{"operands must be numbers"}},
		{"Division #2", "false / true", []string{"operands must be numbers"}},
		{"Addition #1", `"foo" + true`, []string{"operands must be two numbers or two strings"}},
		{"Addition #2", "true + false", []string{"operands must be two numbers or two strings"}},
		{"Subtraction #1", "42 - true", []string{"operands must be numbers"}},
		{"Subtraction #2", `"foo" - "bar"`, []string{"operands must be numbers"}},
		{"Less #1", `"foo" < false`, []string{"operands must be numbers"}},
		{"Less #2", "true < 2", []string{"operands must be numbers"}},
		{"Less #3", `("foo" + "bar") < 42`, []string{"operands must be numbers"}},
		{"Less Or Equal #1", `"foo" <= false`, []string{"operands must be numbers"}},
		{"Less Or Equal #2", "true <= true", []string{"operands must be numbers"}},
		{"Less Or Equal #3", `("foo" + "bar") <= 42`, []string{"operands must be numbers"}},
		{"Greater #1", "false > true", []string{"operands must be numbers"}},
		{"Greater #2", `false > "foo"`, []string{"operands must be numbers"}},
		{"Greater Or Equal #1", "false >= true", []string{"operands must be numbers"}},
		{"Greater Or Equal #2", `"bar" >= "bar"`, []string{"operands must be numbers"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprs := getExpressions(tt.fileContents, t)
			_, errs := EvaluateExpressions(exprs)

			var actual []string
			for _, err := range errs {
				if err == nil {
					continue
				}
				actual = append(actual, err.Error())
			}

			if len(tt.expected) != len(actual) {
				t.Errorf("Evaluation errors length mismatch: Expected %d errors, got %d", len(tt.expected), len(actual))
			} else {
				for index, str := range actual {
					if tt.expected[index] != str {
						t.Errorf("Evaluation errors mismatch on error %d\nExpected: %s\nGot: %s", index+1, tt.expected[index], str)
					}
				}
			}
		})
	}
}
