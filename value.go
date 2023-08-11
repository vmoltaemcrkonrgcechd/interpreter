package main

const (
	numType  = "number"
	boolType = "bool"
	funcType = "fun"
)

type value struct {
	typ  string
	val  any
	body *node
}

func newValue(typ string, val any, body *node) *value {

	return &value{typ: typ, val: val, body: body}
}
