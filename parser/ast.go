package parser

import "goterp/scanner"

type Expr interface{}

type ExprLit struct {
	Token scanner.Token
}

type Stmt interface{}

type StmtPrint struct {
	Value Expr
}
