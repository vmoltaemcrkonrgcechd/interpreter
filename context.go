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

func (ctx *context) edit(name, value string) bool {
	_, ok := ctx.namespace[name]
	if !ok {
		if ctx.parent != nil {
			return ctx.parent.edit(name, value)
		}

		return false
	}

	ctx.namespace[name] = value

	return true
}
