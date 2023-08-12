package main

import (
	"fmt"
	"strconv"
)

var builtin = map[string]func(...any){
	"write": func(args ...any) {
		fmt.Println(args...)
	},
}

type interpreter struct {
}

func newInterpreter() *interpreter {
	return &interpreter{}
}

func (i *interpreter) interpret(root *node, ctx *context) *value {

	switch root.typ {
	case blockType:
		newCtx := newContext(ctx)
		for _, child := range root.children {
			i.interpret(child, newCtx)
		}

	case numberType:
		return newValue(numType, i.parseFloat(root.value), nil)

	case identType:
		val, ok := ctx.find(root.value)
		if !ok {
			panic("неизвестный идентификатор: " + root.value)
		}
		return val

	case addType:
		return newValue(numType, i.interpret(root.children[0], ctx).val.(float64)+i.interpret(root.children[1], ctx).val.(float64), nil)

	case subType:
		return newValue(numType, i.interpret(root.children[0], ctx).val.(float64)-i.interpret(root.children[1], ctx).val.(float64), nil)

	case mulType:
		return newValue(numType, i.interpret(root.children[0], ctx).val.(float64)*i.interpret(root.children[1], ctx).val.(float64), nil)

	case divType:
		return newValue(numType, i.interpret(root.children[0], ctx).val.(float64)/i.interpret(root.children[1], ctx).val.(float64), nil)

	case unarySubType:
		return newValue(numType, -i.interpret(root.children[0], ctx).val.(float64), nil)

	case letType:
		result := 0.0
		if len(root.children) >= 1 {
			result = i.interpret(root.children[0], ctx).val.(float64)
		}
		ctx.namespace[root.value] = newValue(numType, result, nil)
		return nil

	case assignType:
		ok := ctx.edit(root.value, i.interpret(root.children[0], ctx).val.(float64), numType)
		if !ok {
			panic("неизвестный идентификатор:" + root.value)
		}
		val, _ := ctx.find(root.value)
		return val

	case ifType:
		result := i.interpret(root.children[0], ctx)
		if result.val.(bool) {
			i.interpret(root.children[1], ctx)
		}

	case eqlType:
		result := i.interpret(root.children[0], ctx).val.(float64) == i.interpret(root.children[1], ctx).val.(float64)
		return newValue(boolType, result, nil)

	case funType:
		ctx.namespace[root.value] = newValue(funcType, "", root.children[0], root.children[1])

	case functionCallType:
		fun, ok := ctx.find(root.value)
		if !ok {
			var builtinFun func(...any)
			builtinFun, ok = builtin[root.value]
			if ok {
				var parameters []any

				if len(root.children) == 1 {
					for _, parameter := range root.children[0].children {
						if parameter == nil {
							continue
						}
						parameters = append(parameters, i.interpret(parameter, ctx).val.(float64))
					}
				}

				builtinFun(parameters...)
				return nil
			}

			panic("неизвестный идентификатор:" + root.value)
		}

		tempContext := newContext(ctx)

		for ind, arg := range fun.arguments.children {
			val := 0.0
			if len(root.children) == 1 && len(root.children[0].children) > ind {
				val = i.interpret(root.children[0].children[ind], ctx).val.(float64)
			}

			tempContext.namespace[arg.value] = newValue(numType, val, nil)
		}

		return i.interpret(fun.body, tempContext)
	}

	return nil
}

func (i *interpreter) parseFloat(value string) float64 {
	number, err := strconv.ParseFloat(value, 32)

	if err != nil {
		panic(err)
	}

	return number
}
