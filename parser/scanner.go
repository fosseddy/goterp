package parser

import "strconv"

type scanner struct {
	src []byte
	pos int
	cur int
}

func (s *scanner) lexeme() string {
	return string(s.src[s.pos:s.cur])
}

func (s *scanner) hasSrc() bool {
	return s.cur < len(s.src)
}

func (s *scanner) peek() byte {
	return s.src[s.cur]
}

func (s *scanner) peek2() byte {
	next := s.cur + 1
	if next < len(s.src) {
		return s.src[next]
	}
	return 0
}

func (s *scanner) next(ch byte) bool {
	return s.peek() == ch
}

func (s *scanner) advance() byte {
	ch := s.peek()
	s.cur++
	return ch
}

func (s *scanner) scan(t *token) {
scanAgain:
	if !s.hasSrc() {
		t.kind = tokEof
		return
	}

	s.pos = s.cur
	ch := s.advance()

	switch ch {
	case ' ', '\t', '\n', '\r':
		goto scanAgain
	case '/':
		if s.next('/') {
			for s.hasSrc() {
				if s.advance() == '\n' {
					goto scanAgain
				}
			}
		}
		// TODO(art): division

	case ';':
		t.kind = tokSemiColon
	case '=':
		t.kind = tokEq
	case '+':
		t.kind = tokPlus
	case '(':
		t.kind = tokLParen
	case ')':
		t.kind = tokRParen

	default:
		switch {
		case isDigit(ch):
			for isDigit(s.peek()) {
				s.advance()
			}

			if s.next('.') && isDigit(s.peek2()) {
				s.advance() // art: consume '.'
				for isDigit(s.peek()) {
					s.advance()
				}
			}

			// TODO(art): handle error
			f, _ := strconv.ParseFloat(s.lexeme(), 64)
			t.kind = tokNum
			t.val.asNum = f
		case isChar(ch):
			for isAlphaNum(s.peek()) {
				s.advance()
			}

			t.kind = tokIdent
			t.val.asStr = s.lexeme()
		}
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isAlphaNum(ch byte) bool {
	return isChar(ch) || isDigit(ch)
}
