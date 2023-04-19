package ast

import (
	"errors"
	"fmt"
	"golox/lox/expression"
	"golox/lox/reporter"
	"golox/lox/scanner"
	"golox/lox/statement"
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

func (p *Parser) Parse() ([]statement.Stmt, error) {
	statements := make([]statement.Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			continue
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) declaration() (statement.Stmt, error) {
	if p.match(token.Var) {
		return p.varDeclaration()

	}
	return p.statement()
}

func (p *Parser) varDeclaration() (statement.Stmt, error) {
	name, err := p.consume(token.Identifier, "expect variable name")
	if err != nil {
		return nil, err
	}

	var initializer expression.Expression = nil
	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(token.Semicolon, "expect ';' after variable declaration"); err != nil {
		return nil, err
	}

	return statement.NewVarStmt(name, initializer), nil
}

func (p *Parser) statement() (statement.Stmt, error) {
	if p.match(token.For) {
		return p.forStmt()
	}
	if p.match(token.If) {
		return p.ifStmt()
	}
	if p.match(token.Print) {
		return p.printStmt()
	}
	if p.match(token.While) {
		return p.whileStmt()
	}
	// loops
	if p.match(token.Do) {
		statements, err := p.block(token.End)
		if err != nil {
			return nil, err
		}
		return statement.NewBlockStmt(statements), nil
	}
	// blocks
	if p.match(token.LeftBrace) {
		statements, err := p.block(token.RightBrace)
		if err != nil {
			return nil, err
		}
		return statement.NewBlockStmt(statements), nil
	}
	return p.expressionStmt()
}

func (p *Parser) block(limit token.TokenType) ([]statement.Stmt, error) {
	statements := make([]statement.Stmt, 0)

	for !p.check(limit) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}

	if _, err := p.consume(limit, fmt.Sprintf("expect '%s' after a block", limit)); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) ifStmt() (statement.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.Then, "expect 'then' after if condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	branch := "if"

	var elseBranch statement.Stmt = nil
	if p.match(token.Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, nil
		}
		branch = "else"
	}

	if _, err = p.consume(token.End, fmt.Sprintf("expect 'end' after %s branch body", branch)); err != nil {
		return nil, err
	}

	return statement.NewIfStmt(condition, thenBranch, elseBranch), nil
}

func (p *Parser) printStmt() (statement.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.Semicolon, "expect ';' after a value"); err != nil {
		return nil, err
	}
	return statement.NewPrintStmt(val), nil
}

func (p *Parser) whileStmt() (statement.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.check(token.Do) {
		return nil, errors.New("expect 'do' after while conditition")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return statement.NewWhileStmt(condition, body), nil
}

func (p *Parser) forStmt() (statement.Stmt, error) {
	var (
		initializer          statement.Stmt
		condition, increment expression.Expression
		err                  error
	)

	// initializer
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		if initializer, err = p.varDeclaration(); err != nil {
			return nil, err
		}
	} else {
		if initializer, err = p.expressionStmt(); err != nil {
			return nil, err
		}
	}

	// condition
	if !p.match(token.Semicolon) {
		if condition, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(token.Semicolon, "expect ';' after loop condition"); err != nil {
		return nil, err
	}

	// increment
	if !p.check(token.Do) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
		if !p.check(token.Do) {
			return nil, errors.New("expect 'do' after for loop increment")
		}
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// desugarring in a while loop
	if increment != nil {
		body = statement.NewBlockStmt([]statement.Stmt{
			body,
			statement.NewExpressionStmt(increment),
		})
	}
	if condition == nil {
		condition = expression.NewLiteral(true)
	}
	body = statement.NewWhileStmt(condition, body)
	if initializer != nil {
		body = statement.NewBlockStmt([]statement.Stmt{
			initializer,
			body,
		})
	}

	return body, nil
}

func (p *Parser) expressionStmt() (statement.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.Semicolon, "expect ';' after a value"); err != nil {
		return nil, err
	}
	return statement.NewExpressionStmt(val), nil
}

func (p *Parser) expression() (expression.Expression, error) {
	return p.assignment()
}

func (p *Parser) assignment() (expression.Expression, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(token.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*expression.Variable); ok {
			name := v.Name
			return expression.NewAssign(name, value), nil
		}

		p.reporter.ReportAtLocation(errors.New("invalid assignment target"), "TODO", "", equals.Line, 0, 0)

	}
	return expr, nil
}

func (p *Parser) or() (expression.Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = expression.NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (expression.Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.Or) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expression.NewLogical(expr, operator, right)
	}

	return expr, nil
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

	return p.call()
}

func (p *Parser) call() (expression.Expression, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(calle expression.Expression) (expression.Expression, error) {
	/*
		arguments := make([]expression.Expression, 0)
		if !p.check(token.RightParen) {

		}
		paren, err := p.consume(token.RightParen, "expect ')' after arguments")
	*/
	return nil, nil
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

	if p.match(token.Identifier) {
		return expression.NewVariable(p.previous()), nil
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

		if scanner.IsKeyword(p.peek().Type) {
			return
		}

		p.advance()
	}
}
