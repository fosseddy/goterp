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
		for _, value := range s.Values {
			switch v := evaluate(value).(type) {
			case float64, bool:
				fmt.Print(v)
			case nil:
				fmt.Print("nil")
			case string:
				bs := make([]byte, 0, len(v))
				for i := 0; i < len(v); i++ {
					if v[i] == '\\' {
						if i+1 < len(v) && v[i+1] == 'n' {
							bs = append(bs, '\n')
							i++
							continue
						}
					}
					bs = append(bs, v[i])
				}
				fmt.Print(string(bs))
			default:
				// TODO(art): panic, unhandled type
				panic(fmt.Sprintf("unknown expression evaluation type %T\n", v))
			}
		}
	default:
		// TODO(art): panic, unhandled statement
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
		case scanner.TokenStr:
			return e.Value
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
			xb := checkBool(x)
			return !xb
		default:
			// TODO(art): panic, unhandled unary operator
			panic(fmt.Sprintf("unhandled unary operator %v\n", e.Op))
		}
	case parser.ExprBinary:
		x := evaluate(e.X)
		y := evaluate(e.Y)

		switch e.Op {
		case scanner.TokenPlus:
			if xs, ok := x.(string); ok {
				return xs + checkString(y)
			}
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
			return getEquality(x, y)
		case scanner.TokenBangEq:
			return !getEquality(x, y)
		case scanner.TokenLess:
			xf, yf := checkNumbers(x, y)
			return xf < yf
		case scanner.TokenLessEq:
			xf, yf := checkNumbers(x, y)
			return xf <= yf
		case scanner.TokenGreater:
			xf, yf := checkNumbers(x, y)
			return xf > yf
		case scanner.TokenGreaterEq:
			xf, yf := checkNumbers(x, y)
			return xf >= yf
		case scanner.TokenOr:
			xb, yb := checkBools(x, y)
			return xb || yb
		case scanner.TokenAnd:
			xb, yb := checkBools(x, y)
			return xb && yb
		default:
			// TODO(art): panic, unhandled operator
			panic(fmt.Sprintf("unsupported binary operator %v\n", e.Op))
		}
	default:
		// TODO(art): panic, unhandled expression
		panic(fmt.Sprintf("unknown expression type %T\n", e))
	}
}

func getEquality(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if ab, ok := a.(bool); ok {
		return ab == checkBool(b)
	}

	if as, ok := a.(string); ok {
		return as == checkString(b)
	}

	af, bf := checkNumbers(a, b)
	return af == bf
}

func checkNumber(a interface{}) float64 {
	if f, ok := a.(float64); ok {
		return f
	}

	// TODO(art): something better than panic?
	panic(fmt.Sprintf("value %#v must be number\n", a))
}

func checkNumbers(a, b interface{}) (float64, float64) {
	return checkNumber(a), checkNumber(b)
}

func checkBool(a interface{}) bool {
	if b, ok := a.(bool); ok {
		return b
	}

	// TODO(art): something better than panic?
	panic(fmt.Sprintf("value %#v must be bool\n", a))
}

func checkBools(a, b interface{}) (bool, bool) {
	return checkBool(a), checkBool(b)
}

func checkString(a interface{}) string {
	if s, ok := a.(string); ok {
		return s
	}

	// TODO(art): something better than panic?
	panic(fmt.Sprintf("value %#v must be string\n", a))
}
