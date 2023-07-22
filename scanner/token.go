package scanner

type Token int

const (
	TokenInvalid Token = iota

	TokenNum

	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenLParen
	TokenRParen
	TokenSemicolon

	TokenPrint
	TokenTrue
	TokenFalse

	TokenEof
)

func (t Token) String() string {
	switch t {
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
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"

	case TokenEof:
		return "<End of File>"
	}

	return "<Invalid>"
}

var keywords = map[string]Token{
	"print": TokenPrint,
	"true": TokenTrue,
	"false": TokenFalse,
}

func lookupKeyword(s string) Token {
	if k, ok := keywords[s]; ok {
		return k
	}
	return TokenInvalid
}
