package statement

import (
	"golox/lox/expression"
	"golox/lox/token"
)

type ReturnStmt struct {
	Keyword *token.Token // for reporting location
	Value   expression.Expression
}

func NewReturnStmt(keyword *token.Token, value expression.Expression) *ReturnStmt {
	return &ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}
}

func (rs *ReturnStmt) Stmt() {}
