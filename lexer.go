package main

import (
	"bytes"
	"errors"
	"unicode"
)

// lexer - разбирает входную последовательность байтов на токены.
type lexer struct {
	cursor int
	source []rune
}

func newLexer(source []byte) *lexer {
	return &lexer{source: bytes.Runes(source)}
}

func (l *lexer) parse() ([]*token, error) {
	var (
		tokens []*token
		tok    *token
	)

	for l.cursor < len(l.source) {
		if unicode.IsSpace(l.source[l.cursor]) {
			l.cursor++
			continue
		} else if tok = l.parseNumber(); tok != nil {
		} else if tok = l.parseOperator(); tok != nil {
		} else if tok = l.parseIdent(); tok != nil {
		} else {
			// todo: добавить вывод позиции символа.
			return nil, errors.New("неизвестный символ: " + string(l.source[l.cursor]))
		}

		tokens = append(tokens, tok)
	}

	return tokens, nil
}

func (l *lexer) parseNumber() *token {
	var (
		number      []rune
		isFloat     bool
		startCursor = l.cursor
	)

	for l.cursor < len(l.source) {
		if !isFloat && l.source[l.cursor] == '.' {
			isFloat = true
			number = append(number, l.source[l.cursor])
			l.cursor++
			continue
		} else if unicode.IsDigit(l.source[l.cursor]) {
			number = append(number, l.source[l.cursor])
			l.cursor++
			continue
		} else {
			break
		}
	}

	// если число пустое.
	if number == nil ||
		// если число начинается с нуля и после него идут другие символы кроме точки.
		(number[0] == '0' && len(number) > 1 && number[1] != '.') ||
		// если число начинается или заканчивается точкой.
		(isFloat && (number[0] == '.' || number[len(number)-1] == '.')) {
		l.cursor = startCursor
		return nil
	}

	return newToken(numberType, string(number))
}

func (l *lexer) parseOperator() *token {
	var (
		operator    []rune
		startCursor = l.cursor
	)

label:
	for l.cursor < len(l.source) {
		switch l.source[l.cursor] {
		// символы, которые сами являются операторами и не могут накапливаться.
		case '(', ')', ';', '{', '}', ',':
			// если оператор уже накоплен, то курсор не перемещается и этот символ
			// будет прочитан при следующем запуске этой функции.
			if operator != nil {
				break label
			}
			tok := newToken(operators[string(l.source[l.cursor])], "")
			l.cursor++
			return tok

		// символы, которые сами являются операторами, но также могут накапливаться.
		case '+', '-', '*', '/', '=':
			operator = append(operator, l.source[l.cursor])
			l.cursor++

		default:
			break label
		}
	}

	operatorType, ok := operators[string(operator)]
	// если был накоплен неизвестный оператор,
	// то курсор возвращается в исходное состояние и результатом функции будет nil.
	if !ok {
		l.cursor = startCursor
		return nil
	}

	return newToken(operatorType, "")
}

func (l *lexer) parseIdent() *token {
	var (
		ident       []rune
		startCursor = l.cursor
	)

	for l.cursor < len(l.source) {
		if unicode.IsLetter(l.source[l.cursor]) ||
			l.source[l.cursor] == '_' ||
			(len(ident) != 0 && unicode.IsDigit(l.source[l.cursor])) {
			ident = append(ident, l.source[l.cursor])
			l.cursor++
			continue
		}

		break
	}

	if ident == nil {
		l.cursor = startCursor
		return nil
	}

	typ := identType

	keyword, ok := keywords[string(ident)]
	if ok {
		typ = keyword
	}

	return newToken(typ, string(ident))
}
