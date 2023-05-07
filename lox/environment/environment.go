package environment

import (
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

func (env *Environment) Get(name *token.Token) (interface{}, bool) {
	v, found := env.values[name.Lexeme]
	if found {
		return v, true
	}
	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}
	return nil, false
}

func (env *Environment) GetAt(distance int, name string) interface{} {
	return env.ancestor(distance).values[name]
}

func (env *Environment) ancestor(distance int) *Environment {
	ancestor := env
	for i := 0; i < distance; i++ {
		ancestor = ancestor.enclosing
	}
	return ancestor
}

func (env *Environment) Assign(name *token.Token, value interface{}) bool {
	_, found := env.values[name.Lexeme]
	if found {
		env.values[name.Lexeme] = value
		return true
	}
	if env.enclosing != nil {
		return env.enclosing.Assign(name, value)
	}
	return false
}

func (env *Environment) AssignAt(distance int, name *token.Token, value interface{}) bool {
	env.ancestor(distance).values[name.Lexeme] = value
	return true
}
