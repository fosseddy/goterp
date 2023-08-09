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

type Block struct {
	Values map[string]Value
	Prev   *Block
}

// TODO(art): change error
func (b *Block) Declare(k string, v Value) bool {
	_, ok := b.Values[k]
	if ok {
		return false
	}

	b.Values[k] = v;
	return true
}

// TODO(art): change error
func (b *Block) Assign(k string, v Value) bool {
	_, ok := b.Values[k]
	if !ok {
		return false
	}

	b.Values[k] = v;
	return true
}

func (b *Block) Get(k string) (Value, bool) {
	for it := b; it != nil; it = it.Prev {
		if v, ok := it.Values[k]; ok {
			return v, true
		}
	}
	return nil, false
}

var env = Block{map[string]Value{}, nil}

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
		case scanner.TokenIdent:
			v, ok := env.Get(e.Value.Lit)
			if !ok {
				fmt.Fprintf(os.Stderr, "%s does not exist\n", e.Value.Lit)
				os.Exit(1)
			}
			return v
		default:
			panic("unknown literal kind")
		}
	case parser.ExprBinary:
		x, y := eval(e.X), eval(e.Y)

		switch e.Op {
		case scanner.TokenPlus:
			if xs, ok := x.(string); ok {
				return xs + checkStr(y)
			}
			xf, yf := checkNums(x, y)
			return xf + yf
		case scanner.TokenMinus:
			xf, yf := checkNums(x, y)
			return xf - yf
		case scanner.TokenStar:
			xf, yf := checkNums(x, y)
			return xf * yf
		case scanner.TokenSlash:
			xf, yf := checkNums(x, y)
			return xf / yf
		case scanner.TokenLess:
			xf, yf := checkNums(x, y)
			return xf < yf
		case scanner.TokenGreater:
			xf, yf := checkNums(x, y)
			return xf > yf
		case scanner.TokenLessEq:
			xf, yf := checkNums(x, y)
			return xf <= yf
		case scanner.TokenGreaterEq:
			xf, yf := checkNums(x, y)
			return xf >= yf
		case scanner.TokenEqEq:
			return checkEquality(x, y)
		case scanner.TokenBangEq:
			return !checkEquality(x, y)
		case scanner.TokenAnd:
			xb, yb := checkBools(x, y)
			return xb && yb
		case scanner.TokenOr:
			xb, yb := checkBools(x, y)
			return xb || yb
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
	case parser.StmtLet:
		v := eval(s.Value)
		if ok := env.Declare(s.Name, v); !ok {
			fmt.Fprintf(os.Stderr, "%s already exist\n", s.Name)
			os.Exit(1)
		}
	case parser.StmtAssign:
		v := eval(s.Value)
		if ok := env.Assign(s.Name, v); !ok {
			fmt.Fprintf(os.Stderr, "%s does not exist\n", s.Name)
			os.Exit(1)
		}
	case parser.StmtBlock:
		prev := env
		env = Block{map[string]Value{}, &env}

		for _, s := range s.Stmts {
			execute(s)
		}

		env = prev
	default:
		panic("unknown statement kind")
	}
}

func checkNum(x Value) float64 {
	if xf, ok := x.(float64); ok {
		return xf
	}

	// TODO(art): show file and line
	panic("expected number")
}

func checkNums(x, y Value) (float64, float64) {
	return checkNum(x), checkNum(y)
}

func checkBool(x Value) bool {
	if xb, ok := x.(bool); ok {
		return xb
	}

	// TODO(art): show file and line
	panic("expected bool")
}

func checkBools(x, y Value) (bool, bool) {
	return checkBool(x), checkBool(y)
}

func checkStr(x Value) string {
	if xb, ok := x.(string); ok {
		return xb
	}

	// TODO(art): show file and line
	panic("expected string")
}

func checkEquality(x, y Value) bool {
	if x == nil && y == nil {
		return true
	}

	if x == nil || y == nil {
		return false
	}

	if xs, ok := x.(string); ok {
		return xs == checkStr(y)
	}

	if xb, ok := x.(bool); ok {
		return xb == checkBool(y)
	}

	xf, yf := checkNums(x, y)
	return xf == yf
}
