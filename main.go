package main

import (
	"log"
)

func main() {
	input := `2* -2`

	tokens, err := newLexer([]byte(input)).parse()

	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse())
}
