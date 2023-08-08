package parser

import "goterp/scanner"

type StmtKind int

const (
	StmtPrint StmtKind = iota
)

type Stmt struct {
	Kind StmtKind
	Body interface{}
}

type PrintStmt struct {
	Value Expr
}

type ExprKind int

const (
	ExprLit ExprKind = iota
)

type Expr struct {
	Kind ExprKind
	Body interface{}
}

type LitExpr struct {
	Value scanner.Token
}
