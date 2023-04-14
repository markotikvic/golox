package main

import (
	"fmt"
	"golox/lox/ast"
	"golox/lox/expression"
	"golox/lox/token"
)

func main() {
	printer := ast.NewPrinter()
	//lox := lox.NewLox(os.Args)
	//lox.Exec()

	expr := expression.NewBinary(
		expression.NewUnary(
			token.NewToken(token.Minus, "-", nil, 1),
			expression.NewLiteral(123),
		),
		token.NewToken(token.Star, "*", nil, 1),
		expression.NewGrouping(
			expression.NewLiteral(45.67),
		),
	)

	fmt.Println(printer.Print(expr))
}
