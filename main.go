package main

import (
	"fmt"
	"os"
	"strconv"

	"goterp/parser"
	"goterp/scanner"
)

func main() {
	// TODO(art): handle error
	b, _ := os.ReadFile(os.Args[1])

	p := parser.New(b)

	ss := p.Parse()

	for _, s := range ss {
		switch s := s.(type) {
		case parser.StmtPrint:
			switch e := s.Value.(type) {
			case parser.ExprLit:
				switch e.Kind {
				case scanner.TokenNum:
					// TODO(art): handle error
					f, _ := strconv.ParseFloat(e.Value, 64)
					fmt.Printf("%g\n", f)
				default:
					fmt.Printf("unknown token kind %v\n", e.Kind)
				}
			default:
				fmt.Printf("unknown expression type %T\n", e)
			}
		default:
			fmt.Printf("unknown statement type %T\n", s)
		}
	}
}
