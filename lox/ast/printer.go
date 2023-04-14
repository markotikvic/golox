package ast

import (
	"fmt"
	"golox/lox/expression"
	"strings"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(expr expression.Expression) string {
	var stringRep string
	switch v := expr.(type) {
	case *expression.Binary:
		stringRep = p.printBinaryExpr(v)
	case *expression.Unary:
		stringRep = p.printUnaryExpr(v)
	case *expression.Grouping:
		stringRep = p.printGroupingExpr(v)
	case *expression.Literal:
		stringRep = p.printLiteralExpr(v)
	default:
		fmt.Printf("unknown expression type: %v\n", v)
	}
	return stringRep
}

func (p *Printer) printBinaryExpr(expr *expression.Binary) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *Printer) printUnaryExpr(expr *expression.Unary) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *Printer) printGroupingExpr(expr *expression.Grouping) string {
	return p.parenthesize("group", expr.Expr)
}

func (p *Printer) printLiteralExpr(expr *expression.Literal) string {
	if expr.Value == nil {
		return "null"
	}
	return fmt.Sprintf("%#v", expr.Value)
}

func (p *Printer) parenthesize(name string, exprs ...expression.Expression) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(p.Print(expr))
	}
	sb.WriteString(")")
	return sb.String()
}
