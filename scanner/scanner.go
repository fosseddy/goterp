package scanner

import (
	"fmt"
	"os"
)

type TokenKind int

const (
	TokenErr TokenKind = iota

	TokenNum

	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash

	TokenLParen
	TokenRParen
	TokenSemicolon

	TokenPrint

	TokenEof
)

type Token struct {
	Kind TokenKind
	Lit  string
}

func (tok TokenKind) String() string {
	switch tok {
	case TokenNum:
		return "Number"

	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"

	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenSemicolon:
		return ";"

	case TokenPrint:
		return "print"

	case TokenEof:
		return "Eof"
	}

	return "Error"
}

var keywords = map[string]TokenKind{
	"print": TokenPrint,
}

func lookupKeyword(s string) TokenKind {
	if tok, ok := keywords[s]; ok {
		return tok
	}

	return TokenErr
}

type Scanner struct {
	Src      []byte
	SrcLen   int
	Filepath string
	Pos      int
	Cur      int
	Line     int
	Ch       byte
}

func (s *Scanner) literal() string {
	return string(s.Src[s.Pos:s.Cur])
}

func (s *Scanner) hasSrc() bool {
	return s.Cur < s.SrcLen
}

func (s *Scanner) advance() {
	s.Cur++
	if s.hasSrc() {
		s.Ch = s.Src[s.Cur]
	}
}

func (s *Scanner) next(ch byte) bool {
	next := s.Cur + 1
	return next < s.SrcLen && s.Src[next] == ch
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

func (s *Scanner) Scan(tok *Token) {
	tok.Kind = TokenErr
	tok.Lit = ""
scanAgain:
	s.Pos = s.Cur

	if !s.hasSrc() {
		tok.Kind = TokenEof
		return
	}

	switch s.Ch {
	case ' ', '\t', '\r', '\n':
		if s.Ch == '\n' {
			s.Line++
		}
		s.advance()
		goto scanAgain
	case '/':
		if s.next('/') {
			for s.hasSrc() && s.Ch != '\n' {
				s.advance()
			}
			goto scanAgain
		}
		tok.Kind = TokenSlash
		s.advance()
	case '+':
		tok.Kind = TokenPlus
		s.advance()
	case '-':
		tok.Kind = TokenMinus
		s.advance()
	case '*':
		tok.Kind = TokenStar
		s.advance()
	case '(':
		tok.Kind = TokenLParen
		s.advance()
	case ')':
		tok.Kind = TokenRParen
		s.advance()
	case ';':
		tok.Kind = TokenSemicolon
		s.advance()
	default:
		switch {
		case isdigit(s.Ch):
			for isdigit(s.Ch) {
				s.advance()
			}
			if s.Ch == '.' {
				s.advance()
				for isdigit(s.Ch) {
					s.advance()
				}
			}
			tok.Kind = TokenNum
			tok.Lit = s.literal()
		case ischar(s.Ch):
			for isalnum(s.Ch) {
				s.advance()
			}
			tok.Kind = lookupKeyword(s.literal())
			if tok.Kind == TokenErr {
			//if tok.Kind == TokenIdent {
				//tok.Lit = s.literal()
				panic(s.literal())
			}
		default:
			fmt.Fprintf(os.Stderr, "%s:%d:unexpected character %c\n", s.Filepath, s.Line, s.Ch)
			os.Exit(1)
		}
	}
}

func Make(s *Scanner, filepath string) {
	src, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s.Src = src
	s.SrcLen = len(s.Src)
	s.Filepath = filepath
	s.Pos = 0
	s.Cur = 0
	s.Line = 1

	if s.SrcLen > 0 {
		s.Ch = s.Src[0]
	}
}
