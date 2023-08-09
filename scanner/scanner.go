package scanner

import (
	"fmt"
	"log"
	"os"
)

type Scanner struct {
	Src      []byte
	Filepath string
	Pos      int
	Cur      int
	Line     int
	Ch       byte
	isEof    bool
}

func Make(s *Scanner, filepath string) {
	src, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	s.Src = src
	s.Filepath = filepath
	s.Pos = 0
	s.Cur = 0
	s.Line = 1
	s.isEof = false

	if len(s.Src) > 0 {
		s.Ch = s.Src[0]
	}
}

func (s *Scanner) advance() {
	if s.isEof {
		return
	}

	if s.Cur+1 >= len(s.Src) {
		s.isEof = true
		return
	}

	s.Cur++
	s.Ch = s.Src[s.Cur]
}

func (s *Scanner) next(ch byte) bool {
	next := s.Cur + 1
	return next < len(s.Src) && s.Src[next] == ch
}

func (s *Scanner) makeToken(tok *Token, kind TokenKind, lit string) {
	tok.Kind = kind
	tok.Lit = lit
	tok.Line = s.Line
}

func (s *Scanner) makeIdentToken(tok *Token) {
	lit := s.literal()

	tok.Line = s.Line
	tok.Kind = lookupKeyword(lit)

	if tok.Kind == TokenIdent {
		tok.Lit = lit
	}
}

func (s *Scanner) literal() string {
	return string(s.Src[s.Pos:s.Cur])
}

func (s *Scanner) Scan(tok *Token) {
scanAgain:
	if s.isEof {
		s.makeToken(tok, TokenEof, "")
		return
	}

	s.Pos = s.Cur

	switch s.Ch {
	case ' ', '\t', '\r', '\n':
		if s.Ch == '\n' {
			s.Line++
		}
		s.advance()
		goto scanAgain
	case '+':
		s.makeToken(tok, TokenPlus, "")
		s.advance()
	case '-':
		s.makeToken(tok, TokenMinus, "")
		s.advance()
	case '*':
		s.makeToken(tok, TokenStar, "")
		s.advance()
	case '/':
		if s.next('/') {
			for !s.isEof && s.Ch != '\n' {
				s.advance()
			}
			goto scanAgain
		}
		s.makeToken(tok, TokenSlash, "")
		s.advance()
	case '<':
		if s.next('=') {
			s.advance()
			s.makeToken(tok, TokenLessEq, "")
		} else {
			s.makeToken(tok, TokenLess, "")
		}
		s.advance()
	case '>':
		if s.next('=') {
			s.advance()
			s.makeToken(tok, TokenGreaterEq, "")
		} else {
			s.makeToken(tok, TokenGreater, "")
		}
		s.advance()
	case '=':
		if s.next('=') {
			s.advance()
			s.makeToken(tok, TokenEqEq, "")
		} else {
			s.makeToken(tok, TokenEq, "")
		}
		s.advance()
	case '!':
		if s.next('=') {
			s.advance()
			s.makeToken(tok, TokenBangEq, "")
		} else {
			s.makeToken(tok, TokenBang, "")
		}
		s.advance()
	case '&':
		if s.next('&') {
			s.advance()
			s.makeToken(tok, TokenAnd, "")
			s.advance()
		} else {
			fmt.Fprintf(os.Stderr, "%s:%d unexpected character %c\n", s.Filepath, s.Line, s.Ch)
			os.Exit(1)
		}
	case '|':
		if s.next('|') {
			s.advance()
			s.makeToken(tok, TokenOr, "")
			s.advance()
		} else {
			fmt.Fprintf(os.Stderr, "%s:%d unexpected character %c\n", s.Filepath, s.Line, s.Ch)
			os.Exit(1)
		}
	case ';':
		s.makeToken(tok, TokenSemicolon, "")
		s.advance()
	case '(':
		s.makeToken(tok, TokenLParen, "")
		s.advance()
	case ')':
		s.makeToken(tok, TokenRParen, "")
		s.advance()
	case '{':
		s.makeToken(tok, TokenLBrace, "")
		s.advance()
	case '}':
		s.makeToken(tok, TokenRBrace, "")
		s.advance()
	case '"':
		s.advance()
		for !s.isEof && s.Ch != '"' {
			s.advance()
		}
		if s.isEof {
			fmt.Fprintf(os.Stderr, "%s:%d unterminated string\n", s.Filepath, s.Line)
			os.Exit(1)
		}
		s.Pos++
		s.makeToken(tok, TokenStr, s.literal())
		s.advance()
	default:
		switch {
		case isdigit(s.Ch):
			for isdigit(s.Ch) {
				s.advance()
				if s.Ch == '.' {
					s.advance()
					for isdigit(s.Ch) {
						s.advance()
					}
				}
			}
			s.makeToken(tok, TokenNum, s.literal())
		case ischar(s.Ch):
			for isalnum(s.Ch) {
				s.advance()
			}
			s.makeIdentToken(tok)
		default:
			fmt.Fprintf(os.Stderr, "%s:%d unexpected character %c\n", s.Filepath, s.Line, s.Ch)
			os.Exit(1)
		}
	}
}

func isdigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func ischar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isalnum(ch byte) bool {
	return ischar(ch) || isdigit(ch)
}
