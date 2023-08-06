package main

import (
	"fmt"
	"os"
	"strconv"
	"goterp/parser"
	"goterp/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Provide file to execute\n")
		os.Exit(1)
	}

	var p parser.Parser
	parser.Make(&p, os.Args[1])

	for _, s := range p.Parse() {
		execute(s)
	}
}

func execute(s parser.Stmt) {
	switch s := s.(type) {
	case parser.StmtPrint:
		switch v := eval(s.Value).(type) {
		case float64, bool:
			fmt.Println(v)
		default:
			panic("unknown value type")
		}
	default:
		panic("unknown statement type")
	}
}

type Value interface{}

func eval(e parser.Expr) Value {
	switch e := e.(type) {
	case parser.ExprLit:
		switch e.Token.Kind {
		case scanner.TokenNum:
			v, err := strconv.ParseFloat(e.Token.Lit, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return v
		case scanner.TokenTrue:
			return true
		case scanner.TokenFalse:
			return false
		default:
			panic("unknown literal token kind")
		}
	case parser.ExprBinary:
		x, y := eval(e.X), eval(e.Y)

		// TODO(art): proper type validation
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
		default:
			panic("unknown binary expression operator")
		}
	default:
		panic("unknown expression type")
	}
}
