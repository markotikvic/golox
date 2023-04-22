package main

import (
	"golox/lox"
	"os"
)

func main() {
	lox.NewLox(os.Args).Exec()
}
