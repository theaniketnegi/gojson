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
		case 't', 'f':
			l.Reader.UnreadByte()
			t, err := l.tokenizeBoolean()
			if err != nil {
				return nil, err
			}
			l.Tokens = append(l.Tokens, t)
		case 'n':
			l.Reader.UnreadByte()
			t, err := l.tokenizeNull()
			if err != nil {
				return nil, err
			}
			l.Tokens = append(l.Tokens, t)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.Reader.UnreadByte()
			t, err := l.tokenizeNumber()
			if err != nil {
				return nil, err
			}
			l.Tokens = append(l.Tokens, t)
		case '\t':
			return nil, errors.New("illegal character: tab")
		default:
			return nil, fmt.Errorf("illegal character: %v", string(b))
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

func (l *Lexer) tokenizeNumber() (Token, error) {
	b, err := l.Reader.ReadByte()
	if err != nil {
		return Token{}, errors.New("some error occurred while trying to read number")
	}
	n, err := l.Reader.Peek(1)
	if err != nil {
		return Token{}, errors.New("some error occurred while trying to read number")
	}

	if b == '0' && isDigit(n[0]) {
		return Token{}, errors.New("leading zeroes not allowed")
	}

	n, err = l.Reader.Peek(2)
	if err != nil {
		return Token{}, errors.New("some error occurred while trying to read number")
	}

	if b == '-' && n[0] == '0' && isDigit(n[1]) {
		return Token{}, errors.New("leading zeroes not allowed")
	}
	var val bytes.Buffer
	val.WriteByte(b)

	isDecimal := false
	isExponent := false

	for {
		b, err = l.Reader.ReadByte()
		if err != nil {
			return Token{}, errors.New("some error occurred while trying to read number")
		}
		if isDigit(b) {
			val.WriteByte(b)
		} else if b == '.' && !isDecimal {
			isDecimal = true
			val.WriteByte(b)
		} else if (b == 'e' || b == 'E') && !isExponent {
			isExponent = true
			val.WriteByte(b)

			nextByte, err := l.Reader.ReadByte()

			if err != nil {
				return Token{}, errors.New("invalid exponent notation")
			}
			if nextByte == '+' || nextByte == '-' {
				val.WriteByte(nextByte)
			} else if isDigit(nextByte) {
				l.Reader.UnreadByte()
			} else {
				return Token{}, errors.New("invalid exponent notation")
			}
		} else if b == '+' || b == '-' {
			return Token{}, errors.New("unexpected sign")
		} else {
			l.Reader.UnreadByte()
			break
		}
	}

	lastChar := val.String()[val.Len()-1]
	if lastChar == '+' || lastChar == '-' {
		return Token{}, errors.New("incomplete number")
	}

	return Token{Key: NUMBER, Value: val.String()}, nil
}

func (l *Lexer) tokenizeBoolean() (Token, error) {
	n, err := l.Reader.Peek(4)

	if err != nil {
		return Token{}, errors.New("error while trying to read boolean value")
	}

	if bytes.Equal([]byte("true"), n) {
		l.Reader.Discard(4)
		return Token{Key: TRUE, Value: "true"}, nil
	}

	n, err = l.Reader.Peek(5)
	if err != nil {
		return Token{}, errors.New("error while trying to read boolean value")
	}

	if bytes.Equal([]byte("false"), n) {
		l.Reader.Discard(5)
		return Token{Key: FALSE, Value: "false"}, nil
	}
	return Token{}, errors.New("error while trying to read boolean value")
}

func (l *Lexer) tokenizeNull() (Token, error) {
	n, err := l.Reader.Peek(4)
	if err != nil {
		return Token{}, errors.New("error while trying to read null value")
	}

	if bytes.Equal([]byte("null"), n) {
		l.Reader.Discard(4)
		return Token{Key: NULL, Value: "null"}, nil
	}

	return Token{}, errors.New("error while trying to read null value")
}

func isDigit(digit byte) bool {
	return digit >= '0' && digit <= '9'
}
