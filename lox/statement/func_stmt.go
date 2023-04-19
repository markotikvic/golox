package statement

import "golox/lox/token"

type FunctionStmt struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Stmt
}

func NewFunctionStmt(name *token.Token, params []*token.Token, body []Stmt) *FunctionStmt {
	return &FunctionStmt{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (fs *FunctionStmt) Stmt() {}
