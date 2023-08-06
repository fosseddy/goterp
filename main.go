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
		case float64:
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
		default:
			panic("unknown literal token kind")
		}
	default:
		panic("unknown expression type")
	}
}
