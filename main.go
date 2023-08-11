package main

import (
	"log"
)

func main() {
	input := `
let number = 10;

if 0 == 0 {
	let number = 20;
	number;
};

number;
`

	tokens, err := newLexer([]byte(input)).parse()
	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse(), nil)
}
