package parser

import "fmt"

func Parse(src []byte) {
	s := scanner{src: src}

	toks := make([]token, 0, 30)

	for {
		tok := token{}
		s.scan(&tok)

		toks = append(toks, tok)

		if tok.kind == tokEof {
			break
		}
	}

	fmt.Println(toks)
}
