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
		case nil:
			fmt.Println("nil")
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
			f, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				// TODO(art): report, maybe overflow
				panic(fmt.Sprintf("Invalid number value %v\n", e.Value))
			}
			return f
		case scanner.TokenTrue, scanner.TokenFalse:
			b, err := strconv.ParseBool(e.Value)
			if err != nil {
				// TODO(art): panic, this should never happen
				panic(fmt.Sprintf("Invalid boolean value %v\n", e.Value))
			}
			return b
		case scanner.TokenNil:
			var empty interface{}
			return empty
		default:
			// TODO(art): panic, unhandled token
			panic(fmt.Sprintf("unknown token literal kind %v\n", e.Kind))
		}
	case parser.ExprUnary:
		x := evaluate(e.X)

		switch e.Op {
		case scanner.TokenMinus:
			xf := checkNumber(x)
			return -xf
		case scanner.TokenBang:
			switch x := x.(type) {
			case nil:
				return !false
			case float64:
				return !true
			case bool:
				return !x
			default:
				// TODO(art): panic, unhandled value type
				panic(fmt.Sprintf("invalid value type %T for value %v", x, x))
			}
		default:
			panic(fmt.Sprintf("unsupported unary operator %v\n", e.Op))
		}
	case parser.ExprBinary:
		x := evaluate(e.X)
		y := evaluate(e.Y)

		switch e.Op {
		case scanner.TokenPlus:
			xf, yf := checkNumbers(x, y)
			return xf + yf
		case scanner.TokenMinus:
			xf, yf := checkNumbers(x, y)
			return xf - yf
		case scanner.TokenStar:
			xf, yf := checkNumbers(x, y)
			return xf * yf
		case scanner.TokenSlash:
			xf, yf := checkNumbers(x, y)
			return xf / yf
		case scanner.TokenEqEq:
			if x == nil && y == nil {
				return true
			}
			if x == nil || y == nil {
				return false
			}
			if xb, ok := x.(bool); ok {
				return xb == checkBool(y)
			}
			xf, yf := checkNumbers(x, y)
			return xf == yf
		default:
			// TODO(art): panic, unhandled operator
			panic(fmt.Sprintf("unsupported binary operator %v\n", e.Op))
		}
	default:
		// TODO(art): panic, unhandled expression
		panic(fmt.Sprintf("unknown expression type %T\n", e))
	}
}

func checkNumber(a interface{}) float64 {
	if f, ok := a.(float64); ok {
		return f
	}

	// TODO(art): something better than panic?
	panic(fmt.Sprintf("value %v must be number\n", a))
}

func checkBool(a interface{}) bool {
	if b, ok := a.(bool); ok {
		return b
	}

	// TODO(art): something better than panic?
	panic(fmt.Sprintf("value %v must be bool\n", a))
}

func checkNumbers(a, b interface{}) (float64, float64) {
	return checkNumber(a), checkNumber(b)
}
