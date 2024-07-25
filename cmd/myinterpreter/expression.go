package main

type ExpressionType int

var ExpressionTypeEnum = struct {
	LITERAL 	ExpressionType
	UNARY			ExpressionType
	BINARY		ExpressionType
	GROUPING	ExpressionType
}{
	LITERAL:	0,
	UNARY:		1,
	BINARY:		2,
	GROUPING:	3,
}

type Operator int

var OperatorEnum = struct {
	MINUS 					Operator
	PLUS 						Operator
	STAR 						Operator
	EQUAL 					Operator
	EQUAL_EQUAL 		Operator
	BANG 						Operator
	BANG_EQUAL 			Operator
	LESS 						Operator
	LESS_EQUAL 			Operator
	GREATER 				Operator
	GREATER_EQUAL 	Operator
	SLASH 					Operator
}{
	MINUS:					0,
	PLUS:						1,
	STAR:						2,
	EQUAL:					3,
	EQUAL_EQUAL:		4,
	BANG:						5,
	BANG_EQUAL:			6,
	LESS:						7,
	LESS_EQUAL:			8,
	GREATER:				9,
	GREATER_EQUAL:	10,
	SLASH:					11,
}

type Expression struct {
	expression_type ExpressionType
	operator Operator
	literal interface{}
	children []Expression
}