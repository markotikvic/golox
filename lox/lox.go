package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golox/lox/ast"
	"golox/lox/interpreter"
	reporter "golox/lox/reporter"
	"golox/lox/resolver"
	"golox/lox/scanner"
)

type Lox struct {
	args            []string
	hadError        bool
	hadRuntimeError bool
	scanner         *scanner.Scanner
	interp          *interpreter.Interpreter
	reporter        *reporter.ErrorReporter
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
	}
}

func (lox *Lox) Exec() {
	var err error

	switch len(lox.args) {
	case 1:
		err = lox.RunPrompt()
	case 2:
		err = lox.RunScript(lox.args[1])
	default:
		usage()
		os.Exit(64)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (lox *Lox) RunPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("run input: '%s': %w", line, err)
		}
		lox.Run(line, true)
		lox.hadError = false
	}
	return nil
}

func (lox *Lox) RunScript(script string) error {
	source, err := os.ReadFile(script)
	if err != nil {
		err = fmt.Errorf("run script: %w", err)
		return err
	}
	lox.Run(string(source), false)
	if lox.hadError {
		os.Exit(65)
	}
	if lox.hadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func (lox *Lox) Run(source string, repl bool) {
	lox.scanner.Reset()
	tokens := lox.scanner.ScanTokens(source)
	parser := ast.NewParser(tokens, lox.reporter)
	statements, err := parser.Parse()
	if err != nil {
		lox.hadError = true
		return
	}

	resolver := resolver.New(lox.interp, lox.reporter)
	if err = resolver.Resolve(statements); err != nil {
		lox.hadError = true
		return
	}

	if err = lox.interp.Interpret(statements, repl); err != nil {
		lox.hadRuntimeError = true
		return
	}
	//fmt.Println(ast.NewPrinter().Print(tree))
}

func (lox *Lox) HadError() bool {
	return lox.hadError || lox.hadRuntimeError
}

func usage() {
	fmt.Println("usage: lox [script]")
}
