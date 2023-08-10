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
	root := newNode(blockType, "")

	for p.cursor < len(p.tokens) {
		root.children = append(root.children, p.parseExpression())
	}

	return root
}

func (p *parser) parseExpression() *node {
	root := p.parseMulAndDiv()
	for p.cursor < len(p.tokens) {
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
	for p.cursor < len(p.tokens) {
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
	if p.cursor >= len(p.tokens) {
		return nil
	}

	defer func() { p.cursor++ }()

	return newNode(p.tokens[p.cursor].typ, p.tokens[p.cursor].value)
}
