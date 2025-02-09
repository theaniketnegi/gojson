//go:generate stringer -type=TokenType
package main

import (
	"bufio"
	"fmt"
	"io"
)

type TokenType int

// weird enum
const (
	L_BRACE  TokenType = iota // {
	R_BRACE                   // }
	L_SQUARE                  // [
	R_SQUARE                  // ]
	COMMA                     // ,
	COLON                     // :
	NULL
	FALSE
	TRUE
	NUMBER
	STRING
	EOF
)

type Token struct {
	Key   TokenType
	Value string
}

type Lexer struct {
	Reader *bufio.Reader
	Tokens []Token
}

func InitLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{
		Reader: reader,
		Tokens: []Token{},
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	for {
		b, err := l.Reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				l.Tokens = append(l.Tokens, Token{Key: EOF})
				return l.Tokens, nil
			}
			return nil, err
		}

		switch b {
		case ' ', '\n', '\r':
			continue
		case '[':
			l.Tokens = append(l.Tokens, Token{Key: L_SQUARE, Value: "["})
		case ']':
			l.Tokens = append(l.Tokens, Token{Key: R_SQUARE, Value: "]"})
		case '{':
			l.Tokens = append(l.Tokens, Token{Key: L_BRACE, Value: "{"})
		case '}':
			l.Tokens = append(l.Tokens, Token{Key: R_BRACE, Value: "}"})
		case ',':
			l.Tokens = append(l.Tokens, Token{Key: COMMA, Value: ","})
		case ':':
			l.Tokens = append(l.Tokens, Token{Key: COLON, Value: ":"})
		case '"':
			fmt.Print("Most probably a string will go here.\n")
		case 't':
			fmt.Print("True  value.\n")
		case 'f':
			fmt.Print("False value.\n")
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			fmt.Print("Numeric value.\n")
		default:
			return nil, fmt.Errorf("illegal character: %c", b)
		}
	}
}
