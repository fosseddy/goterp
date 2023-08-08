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
		default:
			panic("unknown operation")
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
		case float64:
			fmt.Println(v)
		default:
			panic("unknown value kind")
		}
	default:
		panic("unknown statement kind")
	}
}
