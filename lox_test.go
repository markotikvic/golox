package main

import (
	"golox/lox"
	"os"
	"testing"
)

func TestLox(t *testing.T) {
	tests := []struct {
		script      string
		expectError bool
	}{
		{"./tests/hello.lox", false},
		{"./tests/if_else.lox", false},
		{"./tests/func.lox", false},
		{"./tests/closure.lox", false},
		{"./tests/scoped_error.lox", true},
	}

	for _, test := range tests {
		f, err := os.ReadFile(test.script)
		if err != nil {
			t.Errorf("unable to read script: %s", test.script)
		}
		vm := lox.NewLox(nil)
		vm.Run(string(f), false)
		if vm.HadError() != test.expectError {
			t.Errorf("%s: expected error: %v, got: %v", test.script, test.expectError, vm.HadError())
		}
	}
}
