package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens     []Token
	currentIdx int
}

func NewParser(reader *bufio.Reader) (*Parser, error) {
	lexer := InitLexer(reader)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	return &Parser{
		tokens:     tokens,
		currentIdx: 0,
	}, nil
}

type JsonValue interface{}

func (p *Parser) Parse() (JsonValue, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.currentIdx < len(p.tokens)-1 {
		return nil, errors.New("unexpected token")
	}

	return value, nil
}

func (p *Parser) parseValue() (JsonValue, error) {
	if p.currentIdx >= len(p.tokens) {
		return nil, errors.New("unexpected end of input")
	}

	token := p.tokens[p.currentIdx]

	switch token.Key {
	case L_BRACE:
		return p.parseObject()
	case L_SQUARE:
		return p.parseArray()
	case STRING:
		p.currentIdx++
		return token.Value, nil
	case NUMBER:
		p.currentIdx++
		value, err := strconv.ParseFloat(token.Value, 64)

		if err != nil {
			return nil, fmt.Errorf("invalid number %v: ", token.Value)
		}
		return value, nil
	case TRUE:
		p.currentIdx++
		return true, nil
	case FALSE:
		p.currentIdx++
		return false, nil
	case NULL:
		p.currentIdx++
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", token.Key)
	}
}

func (p *Parser) parseObject() (JsonValue, error) {
	if p.tokens[p.currentIdx].Key != L_BRACE {
		return nil, errors.New("unexpected token, expected: {")
	}
	p.currentIdx++

	obj := make(map[string]JsonValue)

	if p.tokens[p.currentIdx].Key == R_BRACE {
		p.currentIdx++
		return obj, nil
	}

	for {
		if p.tokens[p.currentIdx].Key != STRING {
			return nil, errors.New("unexpected token, expected: string")
		}

		map_key := p.tokens[p.currentIdx].Value
		p.currentIdx++

		if p.tokens[p.currentIdx].Key != COLON {
			return nil, errors.New("unexpected token, expected: colon")
		}
		p.currentIdx++

		map_value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[map_key] = map_value

		if p.tokens[p.currentIdx].Key == R_BRACE {
			p.currentIdx++
			break
		}

		if p.tokens[p.currentIdx].Key != COMMA {
			return nil, errors.New("unexpected token, expected: comma or }")
		}
		p.currentIdx++
	}
	return obj, nil
}

func (p *Parser) parseArray() (JsonValue, error) {
	if p.tokens[p.currentIdx].Key != L_SQUARE {
		return nil, errors.New("unexpected token, expected: [")
	}
	p.currentIdx++

	arr := []JsonValue{}
	if p.tokens[p.currentIdx].Key == R_SQUARE {
		p.currentIdx++
		return arr, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		if p.tokens[p.currentIdx].Key == R_SQUARE {
			p.currentIdx++
			break
		}

		if p.tokens[p.currentIdx].Key != COMMA {
			return nil, errors.New("unexpected token, expected: comma or ]")
		}
		p.currentIdx++
	}
	return arr, nil
}
