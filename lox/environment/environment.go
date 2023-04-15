package environment

import (
	"fmt"
	"golox/lox/token"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Lookup(name *token.Token) (interface{}, error) {
	v, found := env.values[name.Lexeme]
	if !found {
		return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
	}
	return v, nil
}

func (env *Environment) Assign(name *token.Token, value interface{}) error {
	_, found := env.values[name.Lexeme]
	if !found {
		return fmt.Errorf("undefined variable '%s'", name.Lexeme)
	}
	env.values[name.Lexeme] = value
	return nil
}
