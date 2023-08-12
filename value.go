package main

const (
	numType  = "number"
	strType  = "string"
	boolType = "bool"
	funcType = "fun"
)

type value struct {
	typ       string
	val       any
	body      *node
	arguments *node
}

func newValue(typ string, val any, body *node, arguments ...*node) *value {
	if len(arguments) == 1 {
		return &value{typ: typ, val: val, body: body, arguments: arguments[0]}
	}

	return &value{typ: typ, val: val, body: body}
}
