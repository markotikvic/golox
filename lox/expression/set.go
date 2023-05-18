package expression

import "golox/lox/token"

type Set struct {
	Object Expression
	Name   *token.Token
	Value  Expression
}

func NewSet(obj Expression, name *token.Token, val Expression) *Set {
	return &Set{
		Object: obj,
		Name:   name,
		Value:  val,
	}
}

func (g *Set) Expression() {}
