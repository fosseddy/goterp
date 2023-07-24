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

type ExprCall struct {
	Callee string
	Args   []Expr
}

type Stmt interface{}

type StmtPrint struct {
	Values []Expr
}

type StmtVar struct {
	Name string
	Init Expr

	IsFn   bool
	Params []string
	Body   StmtBlock
}

type StmtAssign struct {
	Name  string
	Value Expr
}

type StmtBlock struct {
	Body []Stmt
}

type StmtIf struct {
	Cond      Expr
	IfBlock   StmtBlock
	ElseBlock Stmt
}

type StmtWhile struct {
	Cond Expr
	Body StmtBlock
}

type StmtExpr struct {
	Value Expr
}
