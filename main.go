package main

import (
	"log"
)

func main() {
	input := `
fun bar(a, b) {
	a;
	b;
};

bar();
`

	tokens, err := newLexer([]byte(input)).parse()
	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse(), nil)
}

/*
args
ret
string
*/
