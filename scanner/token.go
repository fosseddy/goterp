package scanner

type Token int

const (
	TokenInvalid Token = iota

	TokenNum

	TokenPlus
	TokenMinus
	TokenLParen
	TokenRParen
	TokenSemicolon

	TokenPrint

	TokenEof
)

func (t Token) String() string {
	switch t {
	case TokenNum:
		return "Number"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "+"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenSemicolon:
		return ";"
	case TokenPrint:
		return "print"
	case TokenEof:
		return "<End of File>"
	}

	return "<Invalid>"
}

var keywords = map[string]Token{
	"print": TokenPrint,
}

func lookupKeyword(s string) Token {
	if k, ok := keywords[s]; ok {
		return k
	}
	return TokenInvalid
}
