package go_interpreter

import (
	"fmt"
)

type Parser struct {
	tokens []Token
	pos    uint
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokens []Token) (Node, error) {
	p.pos = 0
	p.tokens = tokens

	node, err := p.parseComparison()
	if t := p.curToken(); t.Type != EOF {
		return nil, fmt.Errorf("invalid syntax at position %d (expected EOF, got %d)", t.pos, t.Type)
	}
	return node, err
}

func (p *Parser) next() {
	p.pos++
}

func (p *Parser) nextToken() Token {
	return p.tokens[p.pos+1]
}

func (p *Parser) curToken() Token {
	return p.tokens[p.pos]
}

func (p *Parser) parseComparison() (Node, error) {
	result, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
loop:
	for {
		switch t := p.curToken().Type; t {
		case EQ, LT, GT, LTE, GTE:
			p.next()
			right, err := p.parseAddSub()
			if err != nil {
				return nil, err
			}
			result = &Comparison{
				Left:  result,
				Right: right,
				Op:    t,
			}
		default:
			break loop
		}
	}
	return result, nil
}

func (p *Parser) parseAddSub() (Node, error) {
	result, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}
loop:
	for {
		switch t := p.curToken().Type; t {
		case ADD, SUB:
			p.next()
			right, err := p.parseMultiplication()
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{
				Left:  result,
				Right: right,
				Op:    t,
			}
		default:
			break loop
		}
	}
	return result, nil
}

func (p *Parser) parseMultiplication() (Node, error) {
	res, err := p.parseDivision()
	if err != nil {
		return nil, err
	}
	for p.curToken().Type == MUL {
		p.next()
		right, err := p.parseDivision()
		if err != nil {
			return nil, err
		}
		res = &BinaryExpr{
			Left:  res,
			Right: right,
			Op:    MUL,
		}
	}
	return res, nil
}

func (p *Parser) parseDivision() (Node, error) {
	res, err := p.parseExponentiation()
	if err != nil {
		return nil, err
	}
	for p.curToken().Type == DIV {
		p.next()
		right, err := p.parseExponentiation()
		if err != nil {
			return nil, err
		}
		res = &BinaryExpr{
			Left:  res,
			Right: right,
			Op:    DIV,
		}
	}
	return res, nil
}

func (p *Parser) parseExponentiation() (Node, error) {
	res, err := p.parseHighestPriority()
	if err != nil {
		return nil, err
	}
	for p.curToken().Type == EXP {
		p.next()
		right, err := p.parseHighestPriority()
		if err != nil {
			return nil, err
		}
		res = &BinaryExpr{
			Left:  res,
			Right: right,
			Op:    EXP,
		}
	}
	return res, nil
}

func (p *Parser) parseHighestPriority() (Node, error) {
	token := p.curToken()
	switch token.Type {
	case LPAREN:
		p.next()
		res, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		if t := p.curToken(); t.Type != RPAREN {
			return nil, fmt.Errorf("missing ')', got %d", t.Type)
		}
		p.next()
		return res, nil
	case NUMBER, STRING: // literal
		p.next()
		return &Literal{
			Kind:  token.Type,
			Value: token.Value,
		}, nil
	case IDENT:
		// if the next token is a bracket, then parse the function
		if p.nextToken().Type == LPAREN {
			p.next()
			// токен = LPAREN
			args := make([]Node, 0)
		argsLoop:
			for {
				p.next()
				if p.curToken().Type == RPAREN {
					// если аргументов больше нет
					break
				}
				res, err := p.parseComparison()
				if err != nil {
					return nil, err
				}
				args = append(args, res)
				switch t := p.curToken(); t.Type {
				case RPAREN:
					break argsLoop
				case DELIMITER:
					continue
				default:
					return nil, fmt.Errorf("expected ; or ) at the end of the function, got %d", t.Type)
				}
			}
			p.next()
			return &Function{
				Name: token.Value,
				Args: args,
			}, nil

		} else {
			p.next()
			return &Ident{
				Name: token.Value,
			}, nil
		}
	case ADD, SUB: // unary +-
		p.next()
		res, err := p.parseHighestPriority()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Left: res,
			Op:   token.Type,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token: %d", token.Type)
}
