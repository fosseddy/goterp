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

	TokenSemicolon

	TokenPrint

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

	case TokenSemicolon:
		return ";"

	case TokenPrint:
		return "print"

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
}

func lookupKeyword(s string) TokenKind {
	if kwd, ok := keywords[s]; ok {
		return kwd
	}

	return TokenIdent
}
