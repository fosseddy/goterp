package scanner

type TokenKind int

const (
	TokenErr TokenKind = iota

	TokenIdent
	TokenNum

	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash

	TokenLess
	TokenGreater
	TokenLessEq
	TokenGreaterEq
	TokenEq
	TokenEqEq
	TokenBang
	TokenBangEq

	TokenSemicolon
	TokenLParen
	TokenRParen

	TokenPrint
	TokenTrue
	TokenFalse
	TokenNil

	TokenEof
)

func (kind TokenKind) String() string {
	switch kind {
	case TokenIdent:
		return "identifier"
	case TokenNum:
		return "number"

	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"

	case TokenLess:
		return "<"
	case TokenGreater:
		return ">"
	case TokenLessEq:
		return "<="
	case TokenGreaterEq:
		return ">="
	case TokenEq:
		return "="
	case TokenEqEq:
		return "=="
	case TokenBang:
		return "!"
	case TokenBangEq:
		return "!="

	case TokenSemicolon:
		return ";"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"

	case TokenPrint:
		return "print"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNil:
		return "nil"

	case TokenEof:
		return "end of file"

	default:
		return "invalid token"
	}
}

type Token struct {
	Kind TokenKind
	Lit  string
	Line int
}

var keywords = map[string]TokenKind{
	"print": TokenPrint,
	"true": TokenTrue,
	"false": TokenFalse,
	"nil": TokenNil,
}

func lookupKeyword(s string) TokenKind {
	if kwd, ok := keywords[s]; ok {
		return kwd
	}

	return TokenIdent
}
