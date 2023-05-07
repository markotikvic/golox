package expression

import "golox/lox/token"

type Get struct {
	Object interface{}
	Name   *token.Token
}

func NewGet(obj interface{}, name *token.Token) *Get {
	return &Get{
		Object: obj,
		Name:   name,
	}
}

func (g *Get) Expr() {
}
