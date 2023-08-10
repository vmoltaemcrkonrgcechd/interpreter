package main

import (
	"log"
)

func main() {
	input := ``

	tokens, err := newLexer([]byte(input)).parse()

	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse())
}
