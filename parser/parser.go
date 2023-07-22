package parser

import (
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
}

func (p *Parser) match(kind scanner.Token) bool {
	return p.tok == kind
}

func (p *Parser) statement() Stmt {
	return p.printStmt()
}

func (p *Parser) printStmt() StmtPrint {
	if !p.match(scanner.TokenPrint) {
		// TODO(art): handle error
		panic("printStmt: expect print keyword")
	}

	p.advance()
	e := p.expression()

	if !p.match(scanner.TokenSemicolon) {
		// TODO(art): handle error
		panic("printStmt: expect semicolon")
	}
	p.advance()

	return StmtPrint{e}
}

func (p *Parser) expression() Expr {
	return p.term()
}

func (p *Parser) term() Expr {
	e := p.unary()

	for p.match(scanner.TokenPlus) || p.match(scanner.TokenMinus) {
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
	if p.match(scanner.TokenNum) {
		e := ExprLit{p.tok, p.lit}
		p.advance()
		return e
	}

	if p.match(scanner.TokenLParen) {
		p.advance()
		e := p.expression()
		if !p.match(scanner.TokenRParen) {
			// TODO(art): handle error
			panic("primary: expect )")
		}
		p.advance()
		return e
	}

	// TODO(art): handle error
	panic("primary")
}

func (p *Parser) Parse() []Stmt {
	ss := []Stmt{}

	for !p.match(scanner.TokenEof) {
		s := p.statement()
		ss = append(ss, s)
	}

	return ss
}

func New(src []byte) *Parser {
	p := &Parser{s: &scanner.Scanner{Src: src}}
	p.advance()
	return p
}
