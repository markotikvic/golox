package expression

import "golox/lox/token"

type Get struct {
	Object Expression
	Name   *token.Token
}

func NewGet(obj Expression, name *token.Token) *Get {
	return &Get{
		Object: obj,
		Name:   name,
	}
}

func (g *Get) Expression() {}
