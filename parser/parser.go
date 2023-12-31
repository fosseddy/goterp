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
	if p.next(scanner.TokenIdent, scanner.TokenNum, scanner.TokenTrue, scanner.TokenFalse, scanner.TokenNil, scanner.TokenStr) {
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

func (p *Parser) unary() Expr {
	if p.next(scanner.TokenBang, scanner.TokenMinus) {
		op := p.Tok.Kind
		p.advance()
		return ExprUnary{op, p.unary()}
	}

	return p.primary()
}

func (p *Parser) factor() Expr {
	e := p.unary()

	if p.next(scanner.TokenStar, scanner.TokenSlash) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.unary()}
	}

	return e
}

func (p *Parser) term() Expr {
	e := p.factor()

	if p.next(scanner.TokenPlus, scanner.TokenMinus) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.factor()}
	}

	return e
}

func (p *Parser) comparison() Expr {
	e := p.term()

	if p.next(scanner.TokenLess, scanner.TokenGreater, scanner.TokenLessEq, scanner.TokenGreaterEq) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.term()}
	}

	return e
}

func (p *Parser) equality() Expr {
	e := p.comparison()

	if p.next(scanner.TokenEqEq, scanner.TokenBangEq) {
		op := p.Tok.Kind
		p.advance()
		return ExprBinary{e, op, p.comparison()}
	}

	return e
}

func (p *Parser) logicAnd() Expr {
	e := p.equality()

	if p.next(scanner.TokenAnd) {
		p.advance()
		return ExprBinary{e, scanner.TokenAnd, p.equality()}
	}

	return e
}

func (p *Parser) logicOr() Expr {
	e := p.logicAnd()

	if p.next(scanner.TokenOr) {
		p.advance()
		return ExprBinary{e, scanner.TokenOr, p.logicAnd()}
	}

	return e
}

func (p *Parser) expression() Expr {
	return p.logicOr()
}

func (p *Parser) printStmt() Stmt {
	p.consume(scanner.TokenPrint)
	s := StmtPrint{p.expression()}
	p.consume(scanner.TokenSemicolon)

	return s
}

func (p *Parser) letStmt() Stmt {
	p.consume(scanner.TokenLet)

	name := p.Tok.Lit

	p.consume(scanner.TokenIdent)
	p.consume(scanner.TokenEq)

	e := p.expression()

	p.consume(scanner.TokenSemicolon)

	return StmtLet{name, e}
}

func (p *Parser) blockStmt() Stmt {
	var s StmtBlock

	p.consume(scanner.TokenLBrace)

	for !p.next(scanner.TokenRBrace) {
		s.Stmts = append(s.Stmts, p.statement())
	}

	p.consume(scanner.TokenRBrace)
	return s
}

func (p *Parser) assignStmt() Stmt {
	var s StmtAssign

	s.Name = p.Tok.Lit

	p.consume(scanner.TokenIdent)
	p.consume(scanner.TokenEq)

	s.Value = p.expression()

	p.consume(scanner.TokenSemicolon)
	return s
}

func (p *Parser) whileStmt() Stmt {
	p.consume(scanner.TokenWhile)
	return StmtWhile{p.expression(), p.blockStmt()}
}

func (p *Parser) statement() Stmt {
	if p.next(scanner.TokenPrint) {
		return p.printStmt()
	}

	if p.next(scanner.TokenLet) {
		return p.letStmt()
	}

	if p.next(scanner.TokenWhile) {
		return p.whileStmt()
	}

	if p.next(scanner.TokenIdent) {
		return p.assignStmt()
	}

	if p.next(scanner.TokenLBrace) {
		return p.blockStmt()
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
