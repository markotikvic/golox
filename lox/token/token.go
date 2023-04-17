package token

import "fmt"

type TokenType uint32

const (
	// single character tokens
	LeftParen TokenType = 0 + iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Semicolon
	Minus
	Plus
	Slash
	Star

	// one or two character tokens
	Bang
	BangEqual
	Less
	LessEqual
	Greater
	GreaterEqual
	Equal
	EqualEqual
	DotDot // for string concatination

	// literals
	Identifier
	String
	Number

	// keywords
	And
	Or
	Class
	If
	Then
	End
	Else
	Elif
	Not
	While
	For
	Do
	Func
	Null
	Print
	Return
	Base
	Me
	True
	False
	Var

	EOF
)

var tokenTypeNames = []string{
	"LeftParen",
	"RightParen",
	"LeftBrace",
	"RightBrace",
	"Comma",
	"Dot",
	"Semicolon",
	"Minus",
	"Plus",
	"Slash",
	"Star",
	"Bang",
	"BangEqual",
	"Less",
	"LessEqual",
	"Greater",
	"GreaterEqual",
	"Equal",
	"EqualEqual",
	"DotDot",
	"Identifier",
	"String",
	"Number",
	"And",
	"Or",
	"Class",
	"If",
	"Then",
	"End",
	"Else",
	"Elif",
	"Not",
	"While",
	"For",
	"Do",
	"Func",
	"Null",
	"Print",
	"Return",
	"Base",
	"Me",
	"True",
	"False",
	"Var",
	"EOF",
}

func (tt TokenType) String() string {
	return tokenTypeNames[tt]
}

// TODO Expand with starting and ending columns for better error reporting
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s", t.Type, t.Lexeme) // TODO: + literal
}

func NewToken(kind TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    kind,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}
