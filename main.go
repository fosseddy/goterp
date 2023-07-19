package main

import (
	"os"

	"goterp/parser"
)

func main() {
	// TODO(art): handle error
	b, _ := os.ReadFile(os.Args[1])
	parser.Parse(b)
}
