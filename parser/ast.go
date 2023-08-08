package parser

import "goterp/scanner"

type Stmt interface{}

type StmtPrint struct {
	Value Expr
}

type Expr interface{}

type ExprLit struct {
	Value scanner.Token
}
