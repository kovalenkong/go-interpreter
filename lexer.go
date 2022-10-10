package go_interpreter

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Token struct {
	Type  TokenType
	Value string
	pos   uint
}

type Lexer struct {
	tokenPos uint
	reader   *bufio.Reader
}

func NewLexer() *Lexer {
	return new(Lexer)
}

func (l *Lexer) Lex(reader io.Reader) ([]Token, error) {
	l.reader = bufio.NewReader(reader)
	l.tokenPos = 0

	tokens := make([]Token, 0)
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				tokens = append(tokens, Token{
					Type: EOF,
				})
				return tokens, nil
			}
			return nil, err
		}
		var token Token
		l.tokenPos++
		switch {
		case r == '(':
			token = Token{
				Type: LPAREN,
				pos:  l.tokenPos,
			}
		case r == ')':
			token = Token{
				Type: RPAREN,
				pos:  l.tokenPos,
			}
		case r == '+':
			token = Token{
				Type: ADD,
				pos:  l.tokenPos,
			}
		case r == '-':
			token = Token{
				Type: SUB,
				pos:  l.tokenPos,
			}
		case r == '*':
			token = Token{
				Type: MUL,
				pos:  l.tokenPos,
			}
		case r == '/':
			token = Token{
				Type: DIV,
				pos:  l.tokenPos,
			}
		case r == '^':
			token = Token{
				Type: EXP,
				pos:  l.tokenPos,
			}
		case r == ';':
			token = Token{
				Type: DELIMITER,
				pos:  l.tokenPos,
			}
		case r == '"':
			position := l.tokenPos
			value, err := l.readString()
			if err != nil {
				return nil, err
			}
			token = Token{
				Type:  STRING,
				Value: value,
				pos:   position,
			}
		case r == '=':
			token = Token{
				Type: EQ,
				pos:  l.tokenPos,
			}
		case r == '>' || r == '<':
			position := l.tokenPos
			value, err := l.readComparison()
			if err != nil {
				return nil, err
			}
			switch value {
			case ">":
				token = Token{
					Type: GT,
					pos:  position,
				}
			case "<":
				token = Token{
					Type: LT,
					pos:  position,
				}
			case "=":
				token = Token{
					Type: EQ,
					pos:  position,
				}
			case ">=":
				token = Token{
					Type: GTE,
					pos:  position,
				}
			case "<=":
				token = Token{
					Type: LTE,
					pos:  position,
				}
			default:
				return nil, fmt.Errorf("unknown comparasion: %s", value)
			}
		case unicode.IsNumber(r):
			position := l.tokenPos
			value, err := l.readNumber()
			if err != nil {
				return nil, err
			}
			token = Token{
				Type:  NUMBER,
				Value: value,
				pos:   position,
			}
		case unicode.IsLetter(r):
			position := l.tokenPos
			value, err := l.readIdent()
			if err != nil {
				return nil, err
			}
			token = Token{
				Type:  IDENT,
				Value: value,
				pos:   position,
			}
		case unicode.IsSpace(r):
			continue
		default:
			return nil, fmt.Errorf("unknown token: %s", string(r))
		}
		tokens = append(tokens, token)
	}
}

func (l *Lexer) readString() (string, error) {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		l.tokenPos++
		if err != nil {
			if err == io.EOF {
				return lit, nil
			}
			return "", err
		}
		if r == '"' {
			return lit, nil
		}
		lit += string(r)
	}
}

func (l *Lexer) readNumber() (string, error) {
	var number string
	if err := l.reader.UnreadRune(); err != nil {
		return "", err
	}
	l.tokenPos--
	var hasDot bool
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return number, nil
			}
			return "", err
		}
		l.tokenPos++
		switch {
		case unicode.IsDigit(r):
			number += string(r)
		case r == ',' && !hasDot:
			hasDot = true
			number += string(r)
		default:
			err = l.reader.UnreadRune()
			l.tokenPos--
			return number, err
		}
	}
}

func (l *Lexer) readIdent() (string, error) {
	var lit string
	if err := l.reader.UnreadRune(); err != nil {
		return "", err
	}
	l.tokenPos--
	for {
		r, _, err := l.reader.ReadRune()
		l.tokenPos++
		if err != nil {
			if err == io.EOF {
				return lit, nil
			}
			return "", err
		}
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			lit += string(r)
		default:
			err = l.reader.UnreadRune()
			l.tokenPos--
			return lit, err
		}
	}
}

func (l *Lexer) readComparison() (string, error) {
	var lit string
	if err := l.reader.UnreadRune(); err != nil {
		return "", err
	}
	l.tokenPos--
	for {
		r, _, err := l.reader.ReadRune()
		l.tokenPos++
		if err != nil {
			if err == io.EOF {
				return lit, nil
			}
			return "", err
		}
		switch r {
		case '>', '<', '=':
			lit += string(r)
		default:
			err = l.reader.UnreadRune()
			l.tokenPos--
			return lit, err
		}
	}
}
