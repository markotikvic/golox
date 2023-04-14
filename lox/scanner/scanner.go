package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	reporter "golox/lox/reporter"
	"golox/lox/token"
)

var keywords = map[string]token.TokenType{
	"and":    token.And,
	"or":     token.Or,
	"class":  token.Class,
	"if":     token.If,
	"else":   token.Else,
	"not":    token.Not,
	"while":  token.While,
	"for":    token.For,
	"func":   token.Func,
	"null":   token.Null,
	"print":  token.Print,
	"return": token.Return,
	"base":   token.Base,
	"me":     token.Me,
	"true":   token.True,
	"false":  token.False,
	"var":    token.Var,
}

type Scanner struct {
	source              string
	start               int
	current             int
	line                int
	lastDoubleQuoteLine int
	tokens              []*token.Token
	reporter            *reporter.ErrorReporter
}

func NewScanner(reporter *reporter.ErrorReporter) *Scanner {
	return &Scanner{
		line:     1,
		reporter: reporter,
	}
}

func (s *Scanner) ScanTokens(source string) []*token.Token {
	s.source = source
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) Reset() {
	s.line = 1
	s.source = ""
	s.start = 0
	s.current = 0
	s.lastDoubleQuoteLine = 0
	s.tokens = s.tokens[:0]
	s.reporter = nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addMatchingDoubleTokenOr('.', token.DotDot, token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		s.addMatchingDoubleTokenOr('=', token.BangEqual, token.Bang)
	case '=':
		s.addMatchingDoubleTokenOr('=', token.EqualEqual, token.Equal)
	case '<':
		s.addMatchingDoubleTokenOr('=', token.LessEqual, token.Less)
	case '>':
		s.addMatchingDoubleTokenOr('=', token.GreaterEqual, token.Greater)
	case '/':
		if s.peek() == '/' { // comment
			s.skipComment()
		} else {
			s.addToken(token.Slash)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line += 1
	case '"':
		s.addStringToken()
	default:
		// keywords detection
		if unicode.IsDigit(rune(c)) {
			s.addNumberToken()
		} else if s.isValidIdentifierStart(c) {
			s.addIdentifierToken()
		} else {
			s.reporter.Report("TODO", s.line, s.source, fmt.Errorf("unexpected character: %c", rune(c)))
		}
	}
}

func (s *Scanner) advance() byte {
	nextChar := s.source[s.current]
	s.current += 1
	return nextChar
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	nextChar := s.source[s.current]
	return nextChar
}

func (s *Scanner) peekNth(n int) byte {
	if s.isAtEnd() {
		return 0
	}
	nextChar := s.source[s.current+n]
	return nextChar
}

func (s *Scanner) addToken(toktype token.TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(toktype, text, nil, s.line))
}

func (s *Scanner) addTokenWithValue(toktype token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(toktype, text, literal, s.line))
}

func (s *Scanner) addMatchingDoubleTokenOr(char byte, doubleToken token.TokenType, singleToken token.TokenType) {
	if s.peek() == char {
		s.advance()
		s.addToken(doubleToken)
	} else {
		s.addToken(singleToken)
	}
}

func (s *Scanner) skipComment() {
	for s.peek() == '/' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) addStringToken() {
	s.consumeString()
	if s.isAtEnd() {
		s.reporter.Report("TODO", s.line, s.source, fmt.Errorf("unterminated string litteral, started at line: %d", s.lastDoubleQuoteLine))
		return
	}
	s.advance()
	strlit := s.source[s.start+1 : s.current-1]
	s.addTokenWithValue(token.String, strlit)
}

func (s *Scanner) consumeString() {
	s.lastDoubleQuoteLine = s.line
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}
	if s.peek() == '"' {
		s.lastDoubleQuoteLine = s.line
	}
}

func (s *Scanner) addNumberToken() {
	s.consumeNumber()
	if s.peek() == '.' && unicode.IsDigit(rune(s.peekNth(1))) {
		s.advance()
		s.consumeNumber()
	}
	numStr := s.source[s.start:s.current]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		s.reporter.Report("TODO", s.line, s.source, fmt.Errorf("internal error: can't parse float: %s", numStr))
		return
	}
	s.addTokenWithValue(token.Number, num)
}

func (s *Scanner) consumeNumber() {
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}
}

func (s *Scanner) isValidIdentifierStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (s *Scanner) addIdentifierToken() {
	for next := rune(s.peek()); unicode.IsDigit(next) || unicode.IsLetter(next) || next == '_'; next = rune(s.peek()) {
		s.advance()
	}
	identifier := s.source[s.start:s.current]
	if tokenType, isReserved := keywords[identifier]; isReserved {
		s.addToken(tokenType)
	} else {
		s.addToken(token.Identifier)
	}
}
