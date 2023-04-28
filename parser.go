package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Value interface {
	Type() ValueType
	AsFloat64() float64
	AsString() string
	AsBool() bool
}
type Node interface {
	Eval() (Value, error)
}
type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeString
	ValueTypeBool
)

type IntegerLiteral struct {
	Value int
}

type FloatLiteral struct {
	Value float64
}

type BooleanLiteral struct {
	Value bool
}

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}
type UnaryExpression struct {
	Op   string
	Expr Node
}

func (n *IntegerLiteral) Type() ValueType {
	return ValueTypeInt
}

func (n *IntegerLiteral) AsFloat64() float64 {
	return float64(n.Value)
}

func (n *IntegerLiteral) AsString() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n *IntegerLiteral) AsBool() bool {
	return n.Value != 0
}

func (n *FloatLiteral) Type() ValueType {
	return ValueTypeFloat
}

func (n *FloatLiteral) AsFloat64() float64 {
	return n.Value
}

func (n *FloatLiteral) AsString() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *FloatLiteral) AsBool() bool {
	return n.Value != 0
}

func (n *BooleanLiteral) Type() ValueType {
	return ValueTypeBool
}

func (n *BooleanLiteral) AsFloat64() float64 {
	if n.Value {
		return 1
	}
	return 0
}

func (n *BooleanLiteral) AsString() string {
	return fmt.Sprintf("%t", n.Value)
}

func (n *BooleanLiteral) AsBool() bool {
	return n.Value
}

func (n *IntegerLiteral) Eval() (Value, error) {
	return n, nil
}

func (n *FloatLiteral) Eval() (Value, error) {
	return n, nil
}

func (n *BooleanLiteral) Eval() (Value, error) {
	return n, nil
}

func (n *BinaryExpression) Eval() (Value, error) {
	leftVal, err := n.Left.Eval()
	if err != nil {
		return nil, err
	}

	rightVal, err := n.Right.Eval()
	if err != nil {
		return nil, err
	}

	switch n.Op {
	case "&&", "and":
		return &BooleanLiteral{Value: leftVal.AsBool() && rightVal.AsBool()}, nil
	case "||", "or":
		return &BooleanLiteral{Value: leftVal.AsBool() || rightVal.AsBool()}, nil
	case "+":
		return &FloatLiteral{Value: leftVal.AsFloat64() + rightVal.AsFloat64()}, nil
	case "-":
		return &FloatLiteral{Value: leftVal.AsFloat64() - rightVal.AsFloat64()}, nil
	case "*":
		return &FloatLiteral{Value: leftVal.AsFloat64() * rightVal.AsFloat64()}, nil
	case "/":
		if rightVal.AsFloat64() == 0 {
			return nil, errors.New("division by zero")
		}
		return &FloatLiteral{Value: leftVal.AsFloat64() / rightVal.AsFloat64()}, nil
	case "^":
		return &FloatLiteral{Value: math.Pow(leftVal.AsFloat64(), rightVal.AsFloat64())}, nil
	case "==":
		return &BooleanLiteral{Value: leftVal.AsFloat64() == rightVal.AsFloat64()}, nil
	case "!=":
		return &BooleanLiteral{Value: leftVal.AsFloat64() != rightVal.AsFloat64()}, nil
	case "<":
		return &BooleanLiteral{Value: leftVal.AsFloat64() < rightVal.AsFloat64()}, nil
	case ">":
		return &BooleanLiteral{Value: leftVal.AsFloat64() > rightVal.AsFloat64()}, nil
	case "<=":
		return &BooleanLiteral{Value: leftVal.AsFloat64() <= rightVal.AsFloat64()}, nil
	case ">=":
		return &BooleanLiteral{Value: leftVal.AsFloat64() >= rightVal.AsFloat64()}, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", n.Op)
	}
}

func (n *UnaryExpression) Eval() (Value, error) {
	val, err := n.Expr.Eval()
	if err != nil {
		return nil, err
	}

	if n.Op == "!" {
		return &BooleanLiteral{Value: !val.AsBool()}, nil
	}
	return nil, fmt.Errorf("unknown operator: %s", n.Op)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func valToBool(val float64) bool {
	return val != 0
}

type Parser struct {
	tokenizer *Tokenizer
	current   *Token
}

func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer: tokenizer,
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
