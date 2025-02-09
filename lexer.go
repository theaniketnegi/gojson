//go:generate stringer -type=TokenType
package main

import (
	"bufio"
	"bytes"
	"errors"
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
			l.Reader.UnreadByte()
			t, err := l.tokenizeString()
			if err != nil {
				return nil, err
			}
			l.Tokens = append(l.Tokens, t)
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

func (l *Lexer) tokenizeString() (Token, error) {
	b, err := l.Reader.Peek(2)

	if err != nil {
		return Token{}, errors.New(err.Error())
	}

	if bytes.Equal(b, []byte{'"', '"'}) {
		// skip the next two bytes ("")
		l.Reader.Discard(2)
		return Token{Key: STRING, Value: string(b)}, nil
	}
	var val bytes.Buffer

	l.Reader.ReadByte()

	for {
		b, err := l.Reader.ReadByte()
		if err != nil {
			return Token{}, fmt.Errorf("error reading string: %v", string(b))
		}
		if b == '"' {
			break
		}

		if b == '\\' {
			nextByte, err := l.Reader.ReadByte()
			if err != nil {
				return Token{}, fmt.Errorf("error reading string: %v", string(b))
			}

			switch nextByte {
			case '\\', '/', 'b', '"', 'f', 'n', 'r', 't':
				val.WriteByte('\\')
				val.WriteByte(nextByte)
			case 'u':
				unicodeVal := make([]byte, 4)
				_, err := io.ReadFull(l.Reader, unicodeVal)

				if err != nil {
					return Token{}, err
				}
				val.WriteByte('\\')
				val.WriteByte('u')
				val.Write(unicodeVal)
			}
		} else if b == '\t' {
			return Token{}, errors.New("illegal tab character")
		} else if b < 32 || b > 128 {
			return Token{}, fmt.Errorf("illegal character: %v", string(b))
		} else {
			val.WriteByte(b)
		}
	}
	return Token{Key: STRING, Value: val.String()}, nil
}
