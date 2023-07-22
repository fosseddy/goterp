package parser

import "goterp/scanner"

type Expr interface{}

type ExprLit struct {
	Kind  scanner.Token
	Value string
}

type Stmt interface{}

type StmtPrint struct {
	Value Expr
}

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
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.TokenNum) {
		e := ExprLit{p.tok, p.lit}
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
	p := &Parser{s: scanner.New(src)}
	p.advance()
	return p
}
