package statement

import "golox/lox/expression"

type ExpressionStmt struct {
	Expression expression.Expression
}

func NewExpressionStmt(expr expression.Expression) *ExpressionStmt {
	return &ExpressionStmt{
		Expression: expr,
	}
}

func (es *ExpressionStmt) Stmt() {
}
