package statement

import "golox/lox/expression"

type ExpressionStmt struct {
	Expression expression.Expression
}

func (es *ExpressionStmt) Stmt() {
}
