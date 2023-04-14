package main

import (
	"fmt"
	"golox/lox"
	"golox/lox/ast"
	"golox/lox/expression"
	"golox/lox/token"
	"os"
)

func main() {
	lox.NewLox(os.Args).Exec()
	//astPrinterDemo()
}

func astPrinterDemo() {
	printer := ast.NewPrinter()
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
