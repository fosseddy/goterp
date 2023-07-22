package scanner

type Token int

const (
	TokenInvalid Token = iota

	TokenNum

	TokenSemicolon

	TokenPrint

	TokenEof
)

func (t Token) String() string {
	switch t {
	case TokenNum:
		return "Number"
	case TokenSemicolon:
		return ";"
	case TokenPrint:
		return "Print"
	case TokenEof:
		return "<End of File>"
	}

	return "<Invalid>"
}

var keywords = map[string]Token{
	"print": TokenPrint,
}

func keywordLookup(s string) (Token, bool) {
	k, ok := keywords[s]
	return k, ok
}
