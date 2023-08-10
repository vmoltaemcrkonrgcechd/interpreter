package main

type context struct {
	namespace map[string]string
	parent    *context
}

func newContext(parent *context) *context {
	return &context{parent: parent, namespace: make(map[string]string)}
}

func (ctx *context) find(ident string) (string, bool) {
	value, ok := ctx.namespace[ident]
	if ok {
		return value, true
	}

	if ctx.parent != nil {
		return ctx.parent.find(ident)
	}

	return "", false
}
