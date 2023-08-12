package main

import (
	"log"
)

func main() {
	input := `
let number = 1000;

fun bar(number) {
	write(number);
	if number == 1000 {
		bar(number - 100);
	};
};

bar(number);
`

	tokens, err := newLexer([]byte(input)).parse()
	if err != nil {
		log.Fatal(err)
	}

	newInterpreter().interpret(newParser(tokens).parse(), nil)
}

/*
string
ret
*/
