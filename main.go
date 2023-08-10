package main

import (
	"log"
)

/*
let number;
let number = expression;
number = expression;
*/
func main() {
	input := `let n = 10 * 20; n + 40;`

	tokens, err := newLexer([]byte(input)).parse()
	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse(), nil)
}
