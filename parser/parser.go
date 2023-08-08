package parser

import (
	"os"
	"fmt"
	"log"
	"goterp/scanner"
)

type Parser struct {
	Tok scanner.Token
	S scanner.Scanner
}

func Make(p *Parser, filepath string) {
	scanner.Make(&p.S, filepath)
	p.advance()
}

func (p *Parser) advance() {
	var tok scanner.Token

	p.S.Scan(&tok)
	p.Tok = tok
}

func (p *Parser) consume(kind scanner.TokenKind) {
	if p.Tok.Kind != kind {
		fmt.Fprintf(os.Stderr, "%s:%d:expected %s but got %s\n", p.S.Filepath, p.Tok.Line, kind, p.Tok.Kind)
		os.Exit(1)
	}

	p.advance()
}

func (p *Parser) primary(e *Expr) {
	if p.Tok.Kind == scanner.TokenNum {
		var el LitExpr
		
		el.Value = p.Tok
		e.Kind = ExprLit
		e.Body = el

		p.advance()
		return
	}

	log.Fatalf("%s:%d:unknown primary expression %s", p.S.Filepath, p.Tok.Line, p.Tok.Kind)
}

func (p *Parser) expression(e *Expr) {
	p.primary(e)
}

func (p *Parser) printStmt(s *Stmt) {
	var ps PrintStmt;

	p.advance()
	p.expression(&ps.Value)
	p.consume(scanner.TokenSemicolon)

	s.Kind = StmtPrint
	s.Body = ps
}

func (p *Parser) statement(s *Stmt) {
	if p.Tok.Kind == scanner.TokenPrint {
		p.printStmt(s)
		return
	}

	log.Fatalf("%s:%d:unknown statement %s", p.S.Filepath, p.Tok.Line, p.Tok.Kind)
}

func (p *Parser) Parse() []Stmt {
	ss := make([]Stmt, 0, 64)

	for p.Tok.Kind != scanner.TokenEof {
		var s Stmt

		p.statement(&s)
		ss = append(ss, s)
	}

	return ss
}
