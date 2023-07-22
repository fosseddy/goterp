package parser

import (
	"fmt"
	"goterp/scanner"
)

type Parser struct {
	cur int

	s *scanner.Scanner

	tok scanner.Token
	lit string
}

func (p *Parser) advance() {
	p.tok, p.lit = p.s.Scan()
	if p.tok == scanner.TokenInvalid {
		// TODO(art): no panic
		panic(fmt.Sprintln("Token", p.tok, p.lit))
	}
}

func (p *Parser) match(kind scanner.Token) bool {
	return p.tok == kind
}

func (p *Parser) consume(kind scanner.Token) {
	if p.tok != kind {
		// TODO(art): no panic
		panic(fmt.Sprintln("expected", kind, "but got", p.tok))
	}
	p.advance()
}

func (p *Parser) statement() Stmt {
	return p.printStmt()
}

func (p *Parser) printStmt() StmtPrint {
	p.consume(scanner.TokenPrint)
	e := p.expression()
	p.consume(scanner.TokenSemicolon)
	return StmtPrint{e}
}

func (p *Parser) expression() Expr {
	return p.term()
}

func (p *Parser) term() Expr {
	e := p.factor()

	for p.match(scanner.TokenPlus) || p.match(scanner.TokenMinus) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.factor()}
	}

	return e
}

func (p *Parser) factor() Expr {
	e := p.unary()

	for p.match(scanner.TokenStar) || p.match(scanner.TokenSlash) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.unary()}
	}

	return e
}

func (p *Parser) unary() Expr {
	if p.match(scanner.TokenMinus) {
		op := p.tok
		p.advance()
		return ExprUnary{op, p.unary()}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.TokenNum) || p.match(scanner.TokenTrue) || p.match(scanner.TokenFalse) {
		e := ExprLit{p.tok, p.lit}
		p.advance()
		return e
	}

	if p.match(scanner.TokenLParen) {
		p.advance()
		e := p.expression()
		p.consume(scanner.TokenRParen)
		return e
	}

	panic(fmt.Sprintln("Unexpected token", p.tok))
}

func (p *Parser) Parse() []Stmt {
	var ss []Stmt

	for !p.match(scanner.TokenEof) {
		ss = append(ss, p.statement())
	}

	return ss
}

func New(src []byte) *Parser {
	p := &Parser{s: &scanner.Scanner{Src: src}}
	p.advance()
	return p
}
