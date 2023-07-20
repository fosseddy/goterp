package parser

type expr interface{}

type exprLit struct {
	val tokenValue
}

type exprVar struct {
	name string
}

type exprBinary struct {
	x  expr
	op token
	y  expr
}

type exprCall struct {
	callee expr
	args   []expr
}

type stmt interface{}

type stmtVar struct {
	name token
	init expr
}

type parser struct {
	toks []token
	cur  int
}

func New(src []byte) *parser {
	s := scanner{src: src}
	p := &parser{}
	p.toks = make([]token, 0, 30)

	for {
		tok := token{}
		s.scan(&tok)

		p.toks = append(p.toks, tok)

		if tok.kind == tokEof {
			break
		}
	}

	return p
}

func (p *parser) hasTokens() bool {
	return p.peek().kind != tokEof
}

func (p *parser) advance() token {
	t := p.toks[p.cur]
	p.cur++
	return t
}

func (p *parser) peek() token {
	return p.toks[p.cur]
}

func (p *parser) peek2() token {
	next := p.cur + 1
	if next < len(p.toks) {
		return p.toks[next]
	}
	return token{kind: tokEof}
}

func (p *parser) next(kind tokenKind) bool {
	return p.peek().kind == kind
}

func (p *parser) next2(kind tokenKind) bool {
	return p.peek2().kind == kind
}

func (p *parser) declaration() stmt {
	if p.next(tokIdent) && p.next2(tokEq) {
		var_ := stmtVar{}

		var_.name = p.advance()
		p.advance()
		var_.init = p.expression()

		if p.advance().kind != tokSemiColon {
			// TODO(art): handle errors
			panic("expected ;")
		}

		return var_
	}

	e := p.expression()
	if p.advance().kind != tokSemiColon {
		// TODO(art): handle errors
		panic("expected ;")
	}

	return e
}

func (p *parser) expression() expr {
	return p.term()
}

func (p *parser) term() expr {
	e := p.call()

	for p.next(tokPlus) {
		op := p.advance()
		y := p.call()
		e = exprBinary{e, op, y}
	}

	return e
}

func (p *parser) call() expr {
	e := p.primary()

	if p.next(tokLParen) {
		p.advance()
		args := []expr{p.expression()}
		p.advance()
		e = exprCall{e, args}
	}

	return e
}

func (p *parser) primary() expr {
	if p.next(tokLParen) {
		p.advance()
		e := p.expression()
		if p.advance().kind != tokRParen {
			// TODO(art): error handling
			panic("expected right parentheses")
		}

		return e
	}

	if p.next(tokNum) {
		return exprLit{p.advance().val}
	}

	if p.next(tokIdent) {
		return exprVar{p.advance().val.asStr}
	}

	// TODO(art): handle error
	panic("primary")
}

func (p *parser) Parse() []stmt {
	ss := []stmt{}

	for p.hasTokens() {
		s := p.declaration()
		ss = append(ss, s)
	}

	return ss
}
