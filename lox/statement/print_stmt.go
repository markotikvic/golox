package statement

import "golox/lox/expression"

type PrintStmt struct {
	Expression expression.Expression
}

func NewPrintStmt(val expression.Expression) *PrintStmt {
	return &PrintStmt{
		Expression: val,
	}
}

func (ps *PrintStmt) Stmt() {}
