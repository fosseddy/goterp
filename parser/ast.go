package parser

import "goterp/scanner"

type Expr interface{}

type ExprLit struct {
	Token scanner.Token
}

type ExprBinary struct {
	X Expr
	Op scanner.TokenKind
	Y Expr
}

type Stmt interface{}

type StmtPrint struct {
	Value Expr
}
