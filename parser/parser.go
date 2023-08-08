package parser

import (
	"fmt"
	"goterp/scanner"
	"os"
)

type Parser struct {
	Tok scanner.Token
	S   scanner.Scanner
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
		fmt.Fprintf(os.Stderr, "%s:%d expected %s but got %s\n", p.S.Filepath, p.Tok.Line, kind, p.Tok.Kind)
		os.Exit(1)
	}

	p.advance()
}

func (p *Parser) next(kinds ...scanner.TokenKind) bool {
	for _, k := range kinds {
		if p.Tok.Kind == k {
			return true
		}
	}
	return false
}

func (p *Parser) primary() Expr {
	if p.next(scanner.TokenNum) {
		e := ExprLit{p.Tok}
		p.advance()
		return e
	}

	if p.next(scanner.TokenLParen) {
		p.advance()
		e := p.expression()
		p.consume(scanner.TokenRParen)
		return e
	}

	panic(fmt.Sprintf("%s:%d unknown primary %s", p.S.Filepath, p.Tok.Line, p.Tok.Kind))
}

func (p *Parser) factor() Expr {
	e := p.primary()

	if (p.next(scanner.TokenStar, scanner.TokenSlash)) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.factor()}
	}

	return e
}

func (p *Parser) term() Expr {
	e := p.factor()

	if (p.next(scanner.TokenPlus, scanner.TokenMinus)) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.term()}
	}

	return e
}

func (p *Parser) expression() Expr {
	return p.term()
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

	panic(fmt.Sprintf("%s:%d unknown statement %s", p.S.Filepath, p.Tok.Line, p.Tok.Kind))
}

func (p *Parser) Parse() []Stmt {
	ss := make([]Stmt, 0, 256)

	for p.Tok.Kind != scanner.TokenEof {
		ss = append(ss, p.statement())
	}

	return ss
}
