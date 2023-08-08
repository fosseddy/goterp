package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
	"goterp/scanner"
	"goterp/parser"
)

type ValueKind int

const (
	ValueNum ValueKind = iota
)

type Value struct {
	Kind ValueKind
	Num float64
}

func eval(e parser.Expr, res *Value) {
	switch e.Kind {
	case parser.ExprLit:
		el := e.Body.(parser.LitExpr)
		switch el.Value.Kind {
		case scanner.TokenNum:
			f, err := strconv.ParseFloat(el.Value.Lit, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			res.Kind = ValueNum
			res.Num = f
		default:
			log.Fatal("unknown literal kind")
		}
	default:
		log.Fatal("unknown expression kind")
	}
}

func execute(s parser.Stmt) {
	switch s.Kind {
	case parser.StmtPrint:
		var res Value
		ps := s.Body.(parser.PrintStmt)
		eval(ps.Value, &res)

		switch res.Kind {
		case ValueNum:
			fmt.Printf("%g\n", res.Num)
		default:
			log.Fatal("unknown value kind")
		}
	default:
		log.Fatal("unknown statement kind")
	}
}

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
