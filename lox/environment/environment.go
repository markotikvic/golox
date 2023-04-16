package environment

import (
	"fmt"
	"golox/lox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Lookup(name *token.Token) (interface{}, error) {
	v, found := env.values[name.Lexeme]
	if found {
		return v, nil
	}
	if env.enclosing != nil {
		return env.enclosing.Lookup(name)
	}
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}

func (env *Environment) Assign(name *token.Token, value interface{}) error {
	_, found := env.values[name.Lexeme]
	if found {
		env.values[name.Lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.Assign(name, value)
	}
	return fmt.Errorf("undefined variable '%s'", name.Lexeme)
}
