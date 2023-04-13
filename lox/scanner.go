package lox

import (
	"fmt"
	"strconv"
	"unicode"
)

type Scanner struct {
	source              string
	start               int
	current             int
	line                int
	lastDoubleQuoteLine int
	tokens              []*Token
	reporter            *ErrorReporter
}

func newScanner(reporter *ErrorReporter) *Scanner {
	return &Scanner{
		line:     1,
		reporter: reporter,
	}
}

func (s *Scanner) scanTokens(source string) []*Token {
	s.source = source
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, newToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) reset() {
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
		s.addToken(LeftParen, nil)
	case ')':
		s.addToken(RightParen, nil)
	case '{':
		s.addToken(LeftBrace, nil)
	case '}':
		s.addToken(RightBrace, nil)
	case ',':
		s.addToken(Comma, nil)
	case '.':
		s.addMatchingDoubleTokenOr('.', DotDot, Dot)
	case '-':
		s.addToken(Minus, nil)
	case '+':
		s.addToken(Plus, nil)
	case ';':
		s.addToken(Semicolon, nil)
	case '*':
		s.addToken(Star, nil)
	case '!':
		s.addMatchingDoubleTokenOr('=', BangEqual, Bang)
	case '=':
		s.addMatchingDoubleTokenOr('=', EqualEqual, Equal)
	case '<':
		s.addMatchingDoubleTokenOr('=', LessEqual, Less)
	case '>':
		s.addMatchingDoubleTokenOr('=', GreaterEqual, Greater)
	case '/':
		if s.peek() == '/' { // comment
			s.skipComment()
		} else {
			s.addToken(Slash, nil)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line += 1
	case '"':
		s.addStringToken()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.addNumberToken()
	default:
		// keywords detection
		if s.isIdentifierStart(c) {
			s.addIdentifierToken()
		} else {
			s.reporter.report("TODO", s.line, s.source, fmt.Errorf("unexpected character: %c", rune(c)))
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

func (s *Scanner) addToken(toktype TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, newToken(toktype, text, literal, s.line))
}

func (s *Scanner) addMatchingDoubleTokenOr(char byte, doubleToken TokenType, singleToken TokenType) {
	if s.peek() == char {
		s.advance()
		s.addToken(doubleToken, nil)
	} else {
		s.addToken(singleToken, nil)
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
		s.reporter.report("TODO", s.line, s.source, fmt.Errorf("unterminated string litteral, started at line: %d", s.lastDoubleQuoteLine))
		return
	}
	s.advance()
	strlit := s.source[s.start+1 : s.current-1]
	s.addToken(String, strlit)
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
		s.reporter.report("TODO", s.line, s.source, fmt.Errorf("internal error: can't parse float: %s", numStr))
		return
	}
	s.addToken(Number, num)
}

func (s *Scanner) consumeNumber() {
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}
}

func (s *Scanner) isIdentifierStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (s *Scanner) addIdentifierToken() {
}
