package statement

import "golox/lox/token"

type ClassStmt struct {
	Name    *token.Token
	Methods []*FunctionStmt
}

func NewClassStmt(name *token.Token, methods []*FunctionStmt) *ClassStmt {
	return &ClassStmt{
		Name:    name,
		Methods: methods,
	}
}

func (cs *ClassStmt) Stmt() {
}
