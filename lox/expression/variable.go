package expression

import "golox/lox/token"

type Variable struct {
	Name *token.Token
}

func NewVariable(name *token.Token) *Variable {
	return &Variable{
		Name: name,
	}
}

func (e *Variable) Expression() {}
