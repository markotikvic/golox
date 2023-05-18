package interpreter

import (
	"fmt"
	"golox/lox/token"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (li *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", li.class)
}

func (li *LoxInstance) Get(name *token.Token) (interface{}, error) {
	if v, ok := li.fields[name.Lexeme]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("class '%s' has no property called '%s'", li.class.String(), name.Lexeme)
}

func (li *LoxInstance) Set(name *token.Token, val interface{}) {
	li.fields[name.Lexeme] = val
}
