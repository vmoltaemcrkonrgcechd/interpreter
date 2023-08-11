package main

import (
	"log"
)

func main() {
	input := `
fun bar() {
	1;
};

if 0 == 0 {
	fun bar() {
		2;
	};

	bar();
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
fun
ret
string
*/
