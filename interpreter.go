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

	case stringType:
		return newValue(strType, root.value, nil)

	case identType:
		val, ok := ctx.find(root.value)
		if !ok {
			panic("неизвестный идентификатор: " + root.value)
		}
		return val

	case addType:
		op1 := i.interpret(root.children[0], ctx)
		op2 := i.interpret(root.children[1], ctx)
		switch {
		case op1.typ == strType && op2.typ == strType:
			return newValue(strType, op1.val.(string)+op2.val.(string), nil)
		case op1.typ == numType && op2.typ == numType:
			return newValue(numType, op1.val.(float64)+op2.val.(float64), nil)
		case op1.typ == strType && op2.typ == numType:
			return newValue(numType, op1.val.(string)+fmt.Sprintf("%f", op2.val.(float64)), nil)
		case op1.typ == numType && op2.typ == strType:
			return newValue(numType, fmt.Sprintf("%f", op1.val.(float64))+op2.val.(string), nil)
		}

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
		var (
			result any = 0.0
			typ        = numType
		)
		if len(root.children) >= 1 {
			val := i.interpret(root.children[0], ctx)
			typ = val.typ
			if val.typ == strType {
				result = val.val.(string)
			} else {
				result = val.val.(float64)
			}
		}
		ctx.namespace[root.value] = newValue(typ, result, nil)
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
						parameters = append(parameters, i.interpret(parameter, ctx).val)
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
