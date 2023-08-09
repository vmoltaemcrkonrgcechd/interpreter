package main

// token - единица языка (число, строка, ключевое слово и т. д.).
type token struct {
	typ   int
	value string
}

func newToken(typ int, value string) *token {
	return &token{typ: typ, value: value}
}

const (
	numberType = iota
)
