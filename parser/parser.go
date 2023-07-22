package parser

import (
	"fmt"
	"goterp/scanner"
	"strconv"
)

type Expr interface {
	Eval() interface{}
}

type ExprLit struct {
	Kind  scanner.Token
	Value string
}

func (e ExprLit) Eval() interface{} {
	switch e.Kind {
	case scanner.TokenNum:
		// TODO(art): handle error
		f, _ := strconv.ParseFloat(e.Value, 64)
		return f
	default:
		panic(fmt.Sprintf("unknown token literal kind %v\n", e.Kind))
	}
}

type ExprBinary struct {
	X  Expr
	Op scanner.Token
	Y  Expr
}

func (e ExprBinary) Eval() interface{} {
	x := e.X.Eval()
	y := e.Y.Eval()

	switch e.Op {
	case scanner.TokenPlus:
		xf, ok := x.(float64)
		if !ok {
			panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
		}
		yf, ok := y.(float64)
		if !ok {
			panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
		}
		return xf + yf
	default:
		panic(fmt.Sprintf("unsuppported operation %v\n", e.Op))
	}
}

type Stmt interface {
	Execute()
}

type StmtPrint struct {
	Value Expr
}

func (s StmtPrint) Execute() {
	v := s.Value.Eval()
	switch v := v.(type) {
	case float64:
		fmt.Println(v)
	default:
		// TODO(art): this should not be exception, but normal flow of interpreter
		panic(fmt.Sprintf("unknown expression evaluation type %T\n", v))
	}
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
	return p.term()
}

func (p *Parser) term() Expr {
	e := p.primary()

	for p.match(scanner.TokenPlus) {
		op := p.tok
		p.advance()
		e = ExprBinary{e, op, p.primary()}
	}

	return e
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
	p := &Parser{s: &scanner.Scanner{Src: src}}
	p.advance()
	return p
}
