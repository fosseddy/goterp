package parser

import (
	"fmt"
	"goterp/scanner"
	"strconv"
)

type Expr interface {
	Eval() interface{}
}

type ExprLit struct {
	Kind  scanner.Token
	Value string
}

func (e ExprLit) Eval() interface{} {
	switch e.Kind {
	case scanner.TokenNum:
		// TODO(art): handle error
		f, _ := strconv.ParseFloat(e.Value, 64)
		return f
	default:
		panic(fmt.Sprintf("unknown token literal kind %v\n", e.Kind))
	}
}

type ExprUnary struct {
	Op scanner.Token
	X  Expr
}

func (e ExprUnary) Eval() interface{} {
	x := e.X.Eval()

	xf, ok := x.(float64)
	if !ok {
		panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
	}

	switch e.Op {
	case scanner.TokenMinus:
		return -xf
	default:
		panic(fmt.Sprintf("unsupported unary operator %v\n", e.Op))
	}
}

type ExprBinary struct {
	X  Expr
	Op scanner.Token
	Y  Expr
}

func (e ExprBinary) Eval() interface{} {
	x := e.X.Eval()
	y := e.Y.Eval()

	xf, ok := x.(float64)
	if !ok {
		panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
	}
	yf, ok := y.(float64)
	if !ok {
		panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
	}

	switch e.Op {
	case scanner.TokenPlus:
		return xf + yf
	case scanner.TokenMinus:
		return xf - yf
	default:
		panic(fmt.Sprintf("unsupported operation %v\n", e.Op))
	}
}

type Stmt interface {
	Execute()
}

type StmtPrint struct {
	Value Expr
}

func (s StmtPrint) Execute() {
	v := s.Value.Eval()
	switch v := v.(type) {
	case float64:
		fmt.Println(v)
	default:
		// TODO(art): this should not be exception, but normal flow of interpreter
		panic(fmt.Sprintf("unknown expression evaluation type %T\n", v))
	}
}
