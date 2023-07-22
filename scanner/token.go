package scanner

type Token int

const (
	TokenInvalid Token = iota

	TokenNum
	TokenStr

	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenComma
	TokenBang
	TokenEqEq
	TokenBangEq
	TokenLess
	TokenLessEq
	TokenGreater
	TokenGreaterEq
	TokenAnd
	TokenOr

	TokenLParen
	TokenRParen
	TokenSemicolon

	TokenPrint
	TokenTrue
	TokenFalse
	TokenNil

	TokenEof
)

func (t Token) String() string {
	switch t {
	case TokenNum:
		return "Number"
	case TokenStr:
		return "String"

	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenComma:
		return ","
	case TokenBang:
		return "!"
	case TokenEqEq:
		return "=="
	case TokenBangEq:
		return "!="
	case TokenLess:
		return "<"
	case TokenLessEq:
		return "<="
	case TokenGreater:
		return ">"
	case TokenGreaterEq:
		return ">="
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"

	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenSemicolon:
		return ";"

	case TokenPrint:
		return "print"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNil:
		return "nil"

	case TokenEof:
		return "<End of File>"
	}

	return "<Invalid>"
}

var keywords = map[string]Token{
	"print": TokenPrint,
	"true": TokenTrue,
	"false": TokenFalse,
	"nil": TokenNil,
}

func lookupKeyword(s string) Token {
	if k, ok := keywords[s]; ok {
		return k
	}
	return TokenInvalid
}
