package main

import (
	"os"

	"github.com/markotikvic/golox/lox"
)

func main() {
	lox := lox.NewLox(os.Args)
	lox.Exec()
}
