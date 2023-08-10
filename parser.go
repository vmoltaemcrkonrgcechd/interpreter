package main

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
		} else if child = p.parseExpression(); child != nil && (p.cursor >= len(p.tokens) ||
			p.tokens[p.cursor].typ == semicolonType) {
		} else {
			panic("не удалось разобрать выражение")
		}

		root.children = append(root.children, child)
	}

	return root
}

func (p *parser) parseExpression() *node {
	root := p.parseMulAndDiv()
	if root == nil {
		return nil
	}
	for p.cursor < len(p.tokens) && p.tokens[p.cursor].typ != semicolonType {
		if p.tokens[p.cursor].typ != addType && p.tokens[p.cursor].typ != subType {
			return root
		}

		typ := p.tokens[p.cursor].typ

		p.cursor++

		right := p.parseMulAndDiv()
		if right == nil {
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
		}
	}

	return nil
}
