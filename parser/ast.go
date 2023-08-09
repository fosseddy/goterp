package parser

import "goterp/scanner"

type Stmt interface{}

type StmtPrint struct {
	Value Expr
}

type StmtLet struct {
	Name  string
	Value Expr
}

type StmtAssign struct {
	Name  string
	Value Expr
}

type StmtBlock struct {
	Stmts []Stmt
}

type Expr interface{}

type ExprLit struct {
	Value scanner.Token
}

type ExprUnary struct {
	Op scanner.TokenKind
	X  Expr
}

type ExprBinary struct {
	X  Expr
	Op scanner.TokenKind
	Y  Expr
}
