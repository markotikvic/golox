package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golox/lox/ast"
	"golox/lox/interpreter"
	reporter "golox/lox/reporter"
	"golox/lox/scanner"
)

type Lox struct {
	args            []string
	hadError        bool
	hadRuntimeError bool
	scanner         *scanner.Scanner
	interp          *interpreter.Interpreter
	reporter        *reporter.ErrorReporter
	REPL            bool
}

func NewLox(args []string) *Lox {
	reporter := reporter.NewErrorReporter()
	return &Lox{
		args:            args,
		hadError:        false,
		hadRuntimeError: false,
		scanner:         scanner.NewScanner(reporter),
		interp:          interpreter.NewInterpreter(reporter),
		reporter:        reporter,
		REPL:            true,
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
	lox.REPL = true
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
	lox.REPL = false
	source, err := os.ReadFile(script)
	if err != nil {
		err = fmt.Errorf("read file: %w", err)
		return err
	}
	lox.run(string(source))
	if lox.hadError {
		os.Exit(65)
	}
	if lox.hadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func (lox *Lox) run(source string) {
	lox.scanner.Reset()
	tokens := lox.scanner.ScanTokens(source)
	parser := ast.NewParser(tokens, lox.reporter)
	tree, err := parser.Parse()
	if err != nil {
		lox.hadError = true
		return
	}
	lox.interp.Interpret(tree)
	//fmt.Println(ast.NewPrinter().Print(tree))
}

func usage() {
	fmt.Println("usage: lox [script]")
}
