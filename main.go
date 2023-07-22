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
	for _, s := range parser.New(b).Parse() {
		execute(s)
	}
}

func execute(stmt parser.Stmt) {
	switch s := stmt.(type) {
	case parser.StmtPrint:
		switch v := evaluate(s.Value).(type) {
		case float64, bool:
			fmt.Println(v)
		default:
			// TODO(art): this should not be exception, but normal flow of interpreter
			panic(fmt.Sprintf("unknown expression evaluation type %T\n", v))
		}

	default:
		// TODO(art): no panic
		panic(fmt.Sprintf("unknown statement type %T\n", s))
	}
}

func evaluate(expr parser.Expr) interface{} {
	switch e := expr.(type) {
	case parser.ExprLit:
		switch e.Kind {
		case scanner.TokenNum:
			// TODO(art): no panic
			f, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				panic(fmt.Sprintf("Invalid number value %v\n", e.Value))
			}
			return f
		case scanner.TokenTrue, scanner.TokenFalse:
			b, err := strconv.ParseBool(e.Value)
			if err != nil {
				panic(fmt.Sprintf("Invalid boolean value %v\n", e.Value))
			}
			return b
		default:
			// TODO(art): no panic
			panic(fmt.Sprintf("unknown token literal kind %v\n", e.Kind))
		}
	case parser.ExprUnary:
		x := evaluate(e.X)

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
	case parser.ExprBinary:
		x := evaluate(e.X)
		y := evaluate(e.Y)

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
		case scanner.TokenStar:
			return xf * yf
		case scanner.TokenSlash:
			return xf / yf
		default:
			panic(fmt.Sprintf("unsupported binary operator %v\n", e.Op))
		}
	default:
		// TODO(art): no panic
		panic(fmt.Sprintf("unknown expression type %T\n", e))
	}
}
