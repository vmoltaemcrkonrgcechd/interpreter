package main

import (
	"log"
)

func main() {
	input := `
write("Hello " + "World!");
`

	tokens, err := newLexer([]byte(input)).parse()
	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse(), nil)
}

/*
ret
*/
