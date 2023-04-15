package expression

import "golox/lox/token"

type Assign struct {
	Name  *token.Token
	Value Expression
}

func NewAssign(name *token.Token, val Expression) *Assign {
	return &Assign{
		Name:  name,
		Value: val,
	}
}

func (e *Assign) Expression() {}
