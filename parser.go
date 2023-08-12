package main

import (
	"strconv"
)

// parser - строит абстрактное синтаксическое дерево из последовательности токенов.
type parser struct {
	cursor int
	tokens []*token
}

func newParser(tokens []*token) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) parse() *node {
	var (
		child *node
		root  = newNode(blockType, "")
	)

	for p.cursor < len(p.tokens) {
		if p.tokens[p.cursor].typ == semicolonType {
			p.cursor++
			continue
		} else if p.tokens[p.cursor].typ == rBraceType {
			return root
		} else if child = p.parseVariableRedefinition(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else if child = p.parseFunctionCall(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else if child = p.parseBranching(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else if child = p.parseFunctionDeclaration(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else if child = p.parseExpression(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else if child = p.parseVariableDeclaration(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else {
			panic("не удалось разобрать выражение: " + strconv.Itoa(p.cursor))
		}

		root.children = append(root.children, child)
	}

	return root
}

func (p *parser) parseExpression() *node {
	var (
		root        = p.parseMulAndDiv()
		startCursor = p.cursor
	)

	if root == nil {
		p.cursor = startCursor
		return nil
	}
	for p.cursor < len(p.tokens) &&
		p.tokens[p.cursor].typ != semicolonType &&
		p.tokens[p.cursor].typ != commaType {
		if p.tokens[p.cursor].typ != addType && p.tokens[p.cursor].typ != subType {
			return root
		}

		typ := p.tokens[p.cursor].typ

		p.cursor++

		right := p.parseMulAndDiv()
		if right == nil {
			p.cursor = startCursor
			return nil
		}

		root = newNode(typ, "", root, right)
	}

	return root
}

func (p *parser) parseMulAndDiv() *node {
	root := p.parseLiteral()
	if root == nil {
		return nil
	}
	for p.cursor < len(p.tokens) && p.tokens[p.cursor].typ != semicolonType {
		if p.tokens[p.cursor].typ != mulType && p.tokens[p.cursor].typ != divType {
			return root
		}

		typ := p.tokens[p.cursor].typ

		p.cursor++

		right := p.parseLiteral()
		if right == nil {
			return nil
		}

		root = newNode(typ, "", root, right)
	}

	return root
}

func (p *parser) parseLiteral() *node {
	if p.cursor < len(p.tokens) {
		switch p.tokens[p.cursor].typ {
		case lParenType:
			p.cursor++
			expression := p.parseExpression()
			if p.tokens[p.cursor].typ != rParenType {
				return nil
			}
			p.cursor++
			return expression
		case numberType:
			defer func() { p.cursor++ }()
			return newNode(numberType, p.tokens[p.cursor].value)
		case identType:
			defer func() { p.cursor++ }()
			return newNode(identType, p.tokens[p.cursor].value)
		case subType:
			p.cursor++
			literal := p.parseLiteral()
			if literal == nil {
				break
			}
			return newNode(unarySubType, "", literal)
		}
	}

	return nil
}

func (p *parser) parseVariableDeclaration() *node {
	startCursor := p.cursor

	if p.cursor+1 < len(p.tokens) &&
		p.tokens[p.cursor].typ == letType &&
		p.tokens[p.cursor+1].typ == identType {
		variableName := p.tokens[p.cursor+1].value
		p.cursor += 2
		if p.cursor < len(p.tokens) && p.tokens[p.cursor].typ != semicolonType {
			if p.tokens[p.cursor].typ == assignType {
				p.cursor++

				right := p.parseExpression()
				if right == nil {
					p.cursor = startCursor
					return nil
				}

				return newNode(letType, variableName, right)
			}

			return nil
		}

		return newNode(letType, variableName)
	}

	return nil
}

func (p *parser) parseVariableRedefinition() *node {
	startCursor := p.cursor

	if p.cursor+1 < len(p.tokens) &&
		p.tokens[p.cursor].typ == identType &&
		p.tokens[p.cursor+1].typ == assignType {
		root := newNode(assignType, p.tokens[p.cursor].value)
		p.cursor += 2
		expression := p.parseExpression()
		if expression == nil {
			p.cursor = startCursor
			return nil
		}

		root.children = append(root.children, expression)
		return root
	}

	return nil
}

func (p *parser) parseCondition() *node {
	startCursor := p.cursor

	left := p.parseExpression()
	if left == nil || p.cursor >= len(p.tokens) || p.tokens[p.cursor].typ != eqlType {
		p.cursor = startCursor
		return nil
	}

	p.cursor++

	right := p.parseExpression()
	if right == nil {
		p.cursor = startCursor
		return nil
	}

	return newNode(eqlType, "", left, right)
}

func (p *parser) parseBranching() *node {
	startCursor := p.cursor

	if p.tokens[p.cursor].typ == ifType {
		p.cursor++
		condition := p.parseCondition()
		if condition == nil || p.cursor >= len(p.tokens) {
			p.cursor = startCursor
			return nil
		}
		if p.cursor >= len(p.tokens) || p.tokens[p.cursor].typ != lBraceType {
			p.cursor = startCursor
			return nil
		}
		p.cursor++
		body := p.parse()
		if body == nil || p.cursor >= len(p.tokens) || p.tokens[p.cursor].typ != rBraceType {
			p.cursor = startCursor
			return nil
		}
		p.cursor++
		return newNode(ifType, "", condition, body)
	}
	return nil
}

func (p *parser) parseFunctionDeclaration() *node {
	startCursor := p.cursor

	if p.cursor+1 < len(p.tokens) &&
		p.tokens[p.cursor].typ == funType &&
		p.tokens[p.cursor+1].typ == identType {

		root := newNode(funType, p.tokens[p.cursor+1].value)
		p.cursor += 2

		if p.cursor+1 >= len(p.tokens) || p.tokens[p.cursor].typ != lParenType {
			p.cursor = startCursor
			return nil
		}

		p.cursor++

		arguments := p.parseArguments()
		if arguments == nil {
			p.cursor = startCursor
			return nil
		}

		p.cursor++

		if p.cursor+1 >= len(p.tokens) || p.tokens[p.cursor].typ != lBraceType {
			p.cursor = startCursor
			return nil
		}

		p.cursor++

		body := p.parse()
		if body == nil || p.cursor >= len(p.tokens) || p.tokens[p.cursor].typ != rBraceType {
			p.cursor = startCursor
			return nil
		}

		p.cursor++

		root.children = append(root.children, body, arguments)
		return root
	}

	return nil
}

func (p *parser) parseFunctionCall() *node {
	startCursor := p.cursor

	if p.tokens[p.cursor].typ == identType && p.tokens[p.cursor+1].typ == lParenType {
		root := newNode(functionCallType, p.tokens[p.cursor].value)
		p.cursor += 2

		parameters := p.parseParameters()
		if parameters == nil {
			p.cursor = startCursor
			return nil
		}

		root.children = append(root.children, parameters)

		return root
	}

	return nil
}

func (p *parser) parseArguments() *node {
	var (
		arguments   = newNode(argumentsType, "")
		startCursor = p.cursor
	)

	for p.cursor < len(p.tokens) && p.tokens[p.cursor].typ != rParenType {
		if p.tokens[p.cursor].typ == commaType {
			p.cursor++
			continue
		}

		if p.tokens[p.cursor].typ != identType {
			p.cursor = startCursor
			return nil
		}

		arguments.children = append(arguments.children, newNode(identType, p.tokens[p.cursor].value))
		p.cursor++
	}

	if p.cursor >= len(p.tokens) {
		p.cursor = startCursor
		return nil
	}

	return arguments
}

func (p *parser) parseParameters() *node {
	var (
		arguments   = newNode(argumentsType, "")
		startCursor = p.cursor
	)

	for p.cursor < len(p.tokens) {
		right := p.parseExpression()
		if right == nil &&
			p.cursor >= len(p.tokens) &&
			p.tokens[p.cursor].typ != commaType &&
			p.tokens[p.cursor].typ != rParenType {
			p.cursor = startCursor
			return nil
		}

		arguments.children = append(arguments.children, right)

		if p.tokens[p.cursor].typ == rParenType {
			p.cursor++
			return arguments
		}

		p.cursor++
	}

	return nil
}
