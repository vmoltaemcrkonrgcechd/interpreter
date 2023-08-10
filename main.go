package main

import (
	"log"
)

func main() {
	input := `10; 20 * (20 + 30*40/(10*5)) / 25`

	tokens, err := newLexer([]byte(input)).parse()

	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse())
}
