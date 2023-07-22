package scanner

type Scanner struct {
	src []byte
	pos int
	cur int
}

func (s *Scanner) lexeme() string {
	return string(s.src[s.pos:s.cur])
}

func (s *Scanner) hasSrc() bool {
	return s.cur < len(s.src)
}

func (s *Scanner) peek() byte {
	return s.src[s.cur]
}

func (s *Scanner) peek2() byte {
	next := s.cur + 1
	if next < len(s.src) {
		return s.src[next]
	}
	return 0
}

func (s *Scanner) next(ch byte) bool {
	return s.peek() == ch
}

func (s *Scanner) advance() byte {
	ch := s.peek()
	s.cur++
	return ch
}

func (s *Scanner) Scan() (Token, string) {
scanAgain:
	if !s.hasSrc() {
		return TokenEof, ""
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
	case ';':
		return TokenSemicolon, ""

	default:
		switch {
		case isDigit(ch):
			for isDigit(s.peek()) {
				s.advance()
			}

			if s.next('.') && isDigit(s.peek2()) {
				s.advance()
				for isDigit(s.peek()) {
					s.advance()
				}
			}

			return TokenNum, s.lexeme()
		case isChar(ch):
			for isAlphaNum(s.peek()) {
				s.advance()
			}

			if kwd, ok := keywordLookup(s.lexeme()); ok {
				return kwd, ""
			}

			return TokenInvalid, ""
		}
	}

	return TokenInvalid, ""
}

func New(src []byte) *Scanner {
	return &Scanner{src: src}
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
