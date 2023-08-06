package parser

import (
	"fmt"
	"goterp/scanner"
)

type Parser struct {
	S scanner.Scanner
	Tok scanner.Token
}

func (p *Parser) advance() {
	p.S.Scan(&p.Tok)
}

func (p *Parser) consume(kind scanner.TokenKind) {
	// TODO(art): handle error
	if p.Tok.Kind != kind {
		panic(fmt.Sprintf("expected %s but got %s\n", kind, p.Tok.Kind))
	}

	p.advance()
}

func (p *Parser) primary() Expr {
	if p.Tok.Kind == scanner.TokenNum {
		e := ExprLit{p.Tok}
		p.advance()
		return e
	}

	panic("unhandled primary")
}

func (p *Parser) expression() Expr {
	return p.primary()
}

func (p *Parser) printStmt() Stmt {
	p.advance()
	s := StmtPrint{p.expression()}
	p.consume(scanner.TokenSemicolon)

	return s
}

func (p *Parser) statement() Stmt {
	if p.Tok.Kind == scanner.TokenPrint {
		return p.printStmt()
	}

	panic("unhandled statement")
}

func (p *Parser) Parse() []Stmt {
	ss := make([]Stmt, 0, 64)

	for p.Tok.Kind != scanner.TokenEof {
		ss = append(ss, p.statement())
	}

	return ss
}

func Make(p *Parser, filepath string) {
	scanner.Make(&p.S, filepath)
	p.advance()
}
