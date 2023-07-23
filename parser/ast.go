package parser

import "goterp/scanner"

type Expr interface{}

type ExprLit struct {
	Kind  scanner.Token
	Value string
}

type ExprUnary struct {
	Op scanner.Token
	X  Expr
}

type ExprBinary struct {
	X  Expr
	Op scanner.Token
	Y  Expr
}

type Stmt interface{}

type StmtPrint struct {
	Values []Expr
}

type StmtVar struct {
	Name string
	Init Expr
}
