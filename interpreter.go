package main

import (
	"fmt"
	"strconv"
)

type interpreter struct {
}

func newInterpreter() *interpreter {
	return &interpreter{}
}

func (i *interpreter) interpret(root *node, ctx *context) string {

	switch root.typ {
	case blockType:
		newCtx := newContext(ctx)
		for _, child := range root.children {
			fmt.Println(i.interpret(child, newCtx))
		}

	case numberType:
		return root.value

	case identType:
		value, ok := ctx.find(root.value)
		if !ok {
			panic("неизвестный идентификатор: " + root.value)
		}
		return value

	case addType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0], ctx))+i.parseFloat(i.interpret(root.children[1], ctx)))

	case subType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0], ctx))-i.parseFloat(i.interpret(root.children[1], ctx)))

	case mulType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0], ctx))*i.parseFloat(i.interpret(root.children[1], ctx)))

	case divType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0], ctx))/i.parseFloat(i.interpret(root.children[1], ctx)))

	case unarySubType:
		return fmt.Sprintf("-%f", i.parseFloat(i.interpret(root.children[0], ctx)))

	case letType:
		result := "0"
		if len(root.children) >= 1 {
			result = i.interpret(root.children[0], ctx)
		}
		ctx.namespace[root.value] = result
		return result

	case assignType:
		ok := ctx.edit(root.value, i.interpret(root.children[0], ctx))
		if !ok {
			panic("неизвестный идентификатор:" + root.value)
		}
		value, _ := ctx.find(root.value)
		return value

	case ifType:
		result := i.interpret(root.children[0], ctx)
		if result == "true" {
			i.interpret(root.children[1], ctx)
		}

	case eqlType:
		result := i.interpret(root.children[0], ctx) == i.interpret(root.children[1], ctx)
		return fmt.Sprintf("%t", result)
	}

	return ""

}

func (i *interpreter) parseFloat(value string) float64 {
	number, err := strconv.ParseFloat(value, 32)

	if err != nil {
		panic(err)
	}

	return number
}
