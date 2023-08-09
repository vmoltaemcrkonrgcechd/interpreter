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
	if len(number) == 0 ||
		// если число начинается с нуля и после него идут другие символы кроме точки.
		(number[0] == '0' && len(number) > 1 && number[1] != '.') ||
		// если число начинается или заканчивается точкой.
		(isFloat && (number[0] == '.' || number[len(number)-1] == '.')) {
		l.cursor = startCursor
		return nil
	}

	return newToken(numberType, string(number))
}
