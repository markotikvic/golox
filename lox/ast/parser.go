package ast

import (
	"fmt"
	"golox/lox/expression"
	"golox/lox/reporter"
	"golox/lox/token"
)

type Parser struct {
	tokens   []*token.Token
	current  int
	reporter *reporter.ErrorReporter
}

func NewParser(tokens []*token.Token, reporter *reporter.ErrorReporter) *Parser {
	return &Parser{
		tokens:   tokens,
		current:  0,
		reporter: reporter,
	}
}

func (p *Parser) Parse() (expression.Expression, error) {
	tree, err := p.expression()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func (p *Parser) expression() (expression.Expression, error) {
	return p.equality()
}

func (p *Parser) equality() (expression.Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (expression.Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (expression.Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (expression.Expression, error) {
	expr, err := p.concatination()
	if err != nil {
		return nil, err
	}

	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right, err := p.concatination()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) concatination() (expression.Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.DotDot) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (expression.Expression, error) {
	for p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expression.NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (expression.Expression, error) {
	if p.match(token.False) {
		return expression.NewLiteral(false), nil
	}
	if p.match(token.True) {
		return expression.NewLiteral(true), nil
	}
	if p.match(token.Null) {
		return expression.NewLiteral(nil), nil
	}

	if p.match(token.Number, token.String) {
		return expression.NewLiteral(p.previous().Literal), nil
	}

	if p.match(token.LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(token.RightParen, "expect ')' after expression"); err != nil {
			return nil, err
		}
		return expression.NewGrouping(expr), nil
	}

	tok := p.peek()
	prev := p.previous()
	if prev == nil {
		err := fmt.Errorf("unknown expression '%s'", tok.Lexeme)
		p.reporter.ReportAtLocation(err, "TODO", "", tok.Line, 0, 0)
		return nil, err
	}

	err := fmt.Errorf("unknown expression '%s' after '%s'", tok.Lexeme, prev.Lexeme)
	p.reporter.ReportAtLocation(err, "TODO", "", tok.Line, 0, 0)
	return nil, err
}

func (p *Parser) consume(limit token.TokenType, errorMsg string) (*token.Token, error) {
	if p.check(limit) {
		return p.advance(), nil
	}

	var err error
	tok := p.peek()
	if tok.Type == token.EOF {
		err = fmt.Errorf("at end: %s", errorMsg)
	} else {
		err = fmt.Errorf("at '%s': %s", tok.Lexeme, errorMsg)
	}

	p.reporter.ReportAtLocation(err, "TODO", "", tok.Line, 0, 0)
	return nil, err
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	if p.current > 0 {
		return p.tokens[p.current-1]
	}
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Semicolon {
			return
		}

		switch p.peek().Type {
		case token.Class, token.For, token.While, token.Func, token.If, token.Print, token.Return, token.Var:
			return
		}

		p.advance()
	}
}
