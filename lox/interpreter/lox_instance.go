package interpreter

import "fmt"

type LoxInstance struct {
	class *LoxClass
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class: class,
	}
}

func (li *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", li.class)
}
