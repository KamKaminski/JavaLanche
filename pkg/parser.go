package javalanche

import (
	"fmt"
	"log"
	"strconv"
)

type Parser struct {
	tokenizer *Tokenizer
	current   *Token
	evaluator *Evaluator
}

func NewParser(tokenizer *Tokenizer, evaluator *Evaluator) *Parser {
	return &Parser{
		tokenizer: tokenizer,
		evaluator: evaluator,
	}
}

func (p *Parser) Parse() (Node, error) {
	p.nextToken()
	return p.parseTopLevel()
}

func (p *Parser) parseTopLevel() (Node, error) {
	return p.parseBooleanExpression()
}

func (p *Parser) parseBooleanExpression() (Node, error) {
	left, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	for p.current.Type == Operator && (p.current.Value == "and" || p.current.Value == "&&" || p.current.Value == "or" || p.current.Value == "||") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAnd() (Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.current.Type == Keyword && (p.current.Value == "and" || p.current.Value == "&&") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseOr() (Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current.Type == Keyword && (p.current.Value == "or" || p.current.Value == "||") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseComparison() (Node, error) {
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	for p.current.Type == Operator && (p.current.Value == "==" || p.current.Value == "!=" || p.current.Value == "<" || p.current.Value == ">" || p.current.Value == "<=" || p.current.Value == ">=") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) nextToken() error {
	token, err := p.tokenizer.NextToken()
	if err != nil {
		return err
	}
	p.current = token
	fmt.Printf("nextToken: current token = %v\n", p.current)
	return nil
}

func (p *Parser) parseExpression() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.current.Type == Operator && (p.current.Value == "+" || p.current.Value == "-") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.current.Type == Operator && (p.current.Value == "*" || p.current.Value == "/") {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseFactor() (Node, error) {
	base, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return p.parseExponent(base)
}

func (p *Parser) parsePrimary() (Node, error) {
	if p.current.Type == Operator && (p.current.Value == "+" || p.current.Value == "-") {
		op := p.current.Value
		p.nextToken()

		factor, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		if op == "-" {
			return &BinaryExpression{
				Op:    "*",
				Left:  &FloatLiteral{Value: -1},
				Right: factor,
			}, nil
		}

		return factor, nil
	}

	if p.current.Type == Keyword && (p.current.Value == "true" || p.current.Value == "false") {
		value := p.current.Value == "true"
		p.nextToken()
		return &BooleanLiteral{Value: value}, nil
	}

	if p.current.Type == Number {
		intVal, intErr := strconv.Atoi(p.current.Value)
		if intErr == nil {
			p.nextToken()
			return &IntegerLiteral{Value: intVal}, nil
		}

		floatVal, floatErr := strconv.ParseFloat(p.current.Value, 64)
		if floatErr != nil {
			return nil, floatErr
		}
		p.nextToken()
		return &FloatLiteral{Value: floatVal}, nil
	}

	if p.current.Type == LeftParen {
		p.nextToken()
		expr, err := p.parseBooleanExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type != RightParen {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.nextToken()
		return expr, nil
	}

	if p.current.Type == Operator && p.current.Value == "!" {
		p.nextToken()
		expr, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{Op: "!", Expr: expr}, nil
	}
	if p.current.Type == String {
		value := p.current.Value
		p.nextToken()
		log.Println("Primary stirng case", value)
		return &StringLiteral{Value: value}, nil
	}

	return nil, fmt.Errorf("unexpected token: %v", p.current.Value)
}

func (p *Parser) parseExponent(base Node) (Node, error) {
	for p.current.Type == Operator && p.current.Value == "^" {
		op := p.current.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		exponent, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		base = &BinaryExpression{Left: base, Op: op, Right: exponent}
	}
	return base, nil
}
