package main

type context struct {
	namespace map[string]*value
	parent    *context
}

func newContext(parent *context) *context {
	return &context{parent: parent, namespace: make(map[string]*value)}
}

func (ctx *context) find(ident string) (*value, bool) {
	val, ok := ctx.namespace[ident]
	if ok {
		return val, true
	}

	if ctx.parent != nil {
		return ctx.parent.find(ident)
	}

	return nil, false
}

func (ctx *context) edit(name string, value any, typ string) bool {
	_, ok := ctx.namespace[name]
	if !ok {
		if ctx.parent != nil {
			return ctx.parent.edit(name, value, typ)
		}

		return false
	}

	ctx.namespace[name] = newValue(typ, value, nil)

	return true
}
