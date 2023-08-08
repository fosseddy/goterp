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
	case ';':
		s.makeToken(tok, TokenSemicolon, "")
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
			fmt.Fprint(os.Stderr, "%s:%d:unexpected character %c\n", s.Filepath, s.Line, s.Ch)
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
