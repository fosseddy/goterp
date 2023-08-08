package main

import (
	"fmt"
	"goterp/parser"
	"goterp/scanner"
	"os"
	"strconv"
)

func main() {
	var p parser.Parser

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Provide file to execute")
		os.Exit(1)
	}

	parser.Make(&p, os.Args[1])

	for _, s := range p.Parse() {
		execute(s)
	}
}

type Value interface{}

func eval(e parser.Expr) Value {
	switch e := e.(type) {
	case parser.ExprLit:
		switch e.Value.Kind {
		case scanner.TokenNum:
			v, err := strconv.ParseFloat(e.Value.Lit, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return v
		case scanner.TokenTrue:
			return true
		case scanner.TokenFalse:
			return false
		case scanner.TokenNil:
			return nil
		case scanner.TokenStr:
			return e.Value.Lit
		default:
			panic("unknown literal kind")
		}
	case parser.ExprBinary:
		x, y := eval(e.X), eval(e.Y)
		xf, yf := x.(float64), y.(float64)

		switch e.Op {
		case scanner.TokenPlus:
			return xf + yf
		case scanner.TokenMinus:
			return xf - yf
		case scanner.TokenStar:
			return xf * yf
		case scanner.TokenSlash:
			return xf / yf
		case scanner.TokenLess:
			return xf < yf
		case scanner.TokenGreater:
			return xf > yf
		case scanner.TokenLessEq:
			return xf <= yf
		case scanner.TokenGreaterEq:
			return xf >= yf
		case scanner.TokenEqEq:
			return xf == yf
		case scanner.TokenBangEq:
			return xf != yf
		default:
			panic("unknown binary operation")
		}
	case parser.ExprUnary:
		x := eval(e.X)
		switch e.Op {
		case scanner.TokenBang:
			xb := x.(bool)
			return !xb
		case scanner.TokenMinus:
			xf := x.(float64)
			return -xf
		default:
			panic("unknown unary operation")
		}
	default:
		panic("unknown expression kind")
	}
}

func execute(s parser.Stmt) {
	switch s := s.(type) {
	case parser.StmtPrint:
		v := eval(s.Value)

		switch v.(type) {
		case float64, bool, string:
			fmt.Println(v)
		case nil:
			fmt.Println("nil")
		default:
			panic("unknown value kind")
		}
	default:
		panic("unknown statement kind")
	}
}
