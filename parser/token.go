package parser

import "fmt"

type tokenKind int

const (
	tokUnknown tokenKind = iota

	tokIdent
	tokNum

	tokEq
	tokSlash
	tokLParen
	tokRParen
	tokSemiColon

	tokPlus

	tokEof
)

type tokenValue struct {
	asNum float64
	asStr string
}

type token struct {
	kind tokenKind
	val  tokenValue
}

func (t token) String() string {
	switch t.kind {
	case tokIdent:
		return t.val.asStr
	case tokNum:
		return fmt.Sprint(t.val.asNum)
	case tokEq:
		return "="
	case tokSlash:
		return "/"
	case tokLParen:
		return "("
	case tokRParen:
		return ")"
	case tokSemiColon:
		return ";"
	case tokPlus:
		return "+"
	case tokEof:
		return "end of file"
	}

	return "unknown token"
}
