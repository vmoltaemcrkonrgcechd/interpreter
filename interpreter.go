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

func (i *interpreter) interpret(root *node) string {

	switch root.typ {
	case blockType:
		for _, child := range root.children {
			fmt.Println(i.interpret(child))
		}

	case numberType:
		return root.value

	case addType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0]))+i.parseFloat(i.interpret(root.children[1])))

	case subType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0]))-i.parseFloat(i.interpret(root.children[1])))

	case mulType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0]))*i.parseFloat(i.interpret(root.children[1])))

	case divType:
		return fmt.Sprintf("%f", i.parseFloat(i.interpret(root.children[0]))/i.parseFloat(i.interpret(root.children[1])))

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
