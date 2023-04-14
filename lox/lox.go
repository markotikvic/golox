package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"

	reporter "golox/lox/reporter"
	"golox/lox/scanner"
)

type Lox struct {
	args     []string
	hadError bool
	scanner  *scanner.Scanner
	reporter *reporter.ErrorReporter
}

func NewLox(args []string) *Lox {
	reporter := reporter.NewErrorReporter()
	return &Lox{
		args:     args,
		hadError: false,
		scanner:  scanner.NewScanner(reporter),
		reporter: reporter,
	}
}

func (lox *Lox) Exec() {
	var err error

	switch len(lox.args) {
	case 1:
		err = lox.runPrompt()
	case 2:
		err = lox.runScript(lox.args[1])
	default:
		usage()
		os.Exit(64)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (lox *Lox) runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read line '%s': %w", line, err)
		}
		lox.run(line)
		lox.hadError = false
	}
	return nil
}

func (lox *Lox) runScript(script string) error {
	source, err := os.ReadFile(script)
	if err != nil {
		err = fmt.Errorf("read file: %w", err)
		return err
	}
	lox.run(string(source))
	if lox.hadError {
		os.Exit(65)
	}
	return nil
}

func (lox *Lox) run(source string) {
	lox.scanner.Reset()
	tokens := lox.scanner.ScanTokens(source)
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

func usage() {
	fmt.Println("usage: lox [script]")
}
