package main

import (
	"fmt"
	"os"

	"goterp/parser"
)

func main() {
	// TODO(art): handle error
	b, _ := os.ReadFile(os.Args[1])
	p := parser.New(b)

	stmt := p.Parse()

	for _, s := range stmt {
		fmt.Printf("%+v\n", s)
	}
}
