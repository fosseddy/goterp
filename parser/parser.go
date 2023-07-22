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

func (p *Parser) match(kinds ...scanner.Token) bool {
	for _, k := range kinds {
		if p.tok == k {
			return true
		}
	}

	return false
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

	es := []Expr{p.expression()}
	for p.match(scanner.TokenComma) {
		p.advance()
		es = append(es, p.expression())
	}

	p.consume(scanner.TokenSemicolon)
	return StmtPrint{es}
}

func (p *Parser) expression() Expr {
	return p.logicOr()
}

func (p *Parser) logicOr() Expr {
	e := p.logicAnd()

	for p.match(scanner.TokenOr) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.logicAnd()}
	}

	return e
}

func (p *Parser) logicAnd() Expr {
	e := p.equality()

	for p.match(scanner.TokenAnd) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.equality()}
	}

	return e
}

func (p *Parser) equality() Expr {
	e := p.comparison()

	for p.match(scanner.TokenEqEq, scanner.TokenBangEq) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.comparison()}
	}

	return e
}

func (p *Parser) comparison() Expr {
	e := p.term()

	for p.match(scanner.TokenLess, scanner.TokenLessEq, scanner.TokenGreater, scanner.TokenGreaterEq) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.term()}
	}

	return e
}

func (p *Parser) term() Expr {
	e := p.factor()

	for p.match(scanner.TokenPlus, scanner.TokenMinus) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.factor()}
	}

	return e
}

func (p *Parser) factor() Expr {
	e := p.unary()

	for p.match(scanner.TokenStar, scanner.TokenSlash) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.unary()}
	}

	return e
}

func (p *Parser) unary() Expr {
	if p.match(scanner.TokenMinus, scanner.TokenBang) {
		op := p.tok
		p.advance()
		return ExprUnary{op, p.unary()}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.TokenNum, scanner.TokenStr, scanner.TokenTrue, scanner.TokenFalse, scanner.TokenNil) {
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

	// TODO(art): report, probably syntax error
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
