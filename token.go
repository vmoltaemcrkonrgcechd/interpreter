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
	blockType = iota

	identType
	numberType

	addType
	subType
	unarySubType
	mulType
	divType
	lParenType
	rParenType
	lBraceType
	rBraceType
	eqlType
	semicolonType

	assignType

	letType
	ifType
	elseType
)

var operators = map[string]int{
	"+":  addType,
	"-":  subType,
	"*":  mulType,
	"/":  divType,
	"(":  lParenType,
	")":  rParenType,
	";":  semicolonType,
	"=":  assignType,
	"{":  lBraceType,
	"}":  rBraceType,
	"==": eqlType,
}

var keywords = map[string]int{
	"let":  letType,
	"if":   ifType,
	"else": elseType,
}
