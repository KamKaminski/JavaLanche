package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
)

var (
	_ fmt.Stringer = (*IntegerLiteral)(nil)
	_ fmt.Stringer = (*FloatLiteral)(nil)
	_ fmt.Stringer = (*StringLiteral)(nil)
	_ fmt.Stringer = (*BooleanLiteral)(nil)
	_ fmt.Stringer = (*UnaryExpression)(nil)
	_ fmt.Stringer = (*BinaryExpression)(nil)
)

var (
	errInvalidType  = errors.New("invalid type")
	errInvalidTypes = errors.New("invalid types")
	errInvalidOp    = errors.New("invalid operator")
	errDivZero      = errors.New("division by zero")
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

type StringLiteral struct {
	Value string
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

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("&%T{%v}", *n, n.Value)
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

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("&%T{%v}", n, n.Value)
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

func (n *BooleanLiteral) String() string {
	if n.Value {
		return "true"
	}
	return "false"
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
	return n.String()
}

func (n *BooleanLiteral) AsBool() bool {
	return n.Value
}

func (n *StringLiteral) String() string {
	return n.Value
}

func (n *StringLiteral) Type() ValueType {
	return ValueTypeString
}

func (n *StringLiteral) AsBool() bool {
	return n.Value != ""
}

func (n *StringLiteral) AsString() string {
	return n.String()
}

func (n *StringLiteral) AsFloat64() float64 {
	v, _ := strconv.ParseFloat(n.Value, 64)
	return v
}

func (n *IntegerLiteral) Eval() (Value, error) {
	log.Println("Integer.Eval", n.Value)
	return n, nil
}

func (n *FloatLiteral) Eval() (Value, error) {
	log.Println("Float.Eval", n.Value)
	return n, nil
}

func (n *BooleanLiteral) Eval() (Value, error) {
	log.Println("Bool.Eval", n.Value)
	return n, nil
}

func (n *BinaryExpression) String() string {
	return fmt.Sprintf("&%T{%s, %q, %s}", *n, n.Left, n.Op, n.Right)
}

func (n *BinaryExpression) Eval() (Value, error) {
	log.Println("Binary.Eval", n)

	leftVal, err := n.Left.Eval()
	if err != nil {
		return nil, err
	}

	rightVal, err := n.Right.Eval()
	if err != nil {
		return nil, err
	}

	switch leftVal.Type() {
	case ValueTypeBool:
		return evalBinaryBool(n.Op, leftVal.AsBool(), rightVal)
	case ValueTypeString:
		return evalBinaryString(n.Op, leftVal.AsString(), rightVal)
	default:
		return evalBinaryFloat(n.Op, leftVal.AsFloat64(), rightVal)
	}
}

func evalBinaryBool(op string, left bool, rightVal Value) (Value, error) {
	var right, result bool

	switch {
	case rightVal.Type() == ValueTypeBool:
		right = rightVal.AsBool()
	case op == "==", op == "!=":
		return &BooleanLiteral{Value: false}, nil
	default:
		return nil, errInvalidTypes
	}

	switch op {
	case "&&", "and":
		// AND
		result = left && right
	case "||", "or":
		// OR
		result = left || right
	case "==":
		// EQ
		result = left == right
	case "!=":
		// NE
		result = left != right
	case "^":
		// XOR
		result = (left && !right) || (!left && right)
	default:
		return nil, errInvalidOp
	}

	return &BooleanLiteral{Value: result}, nil
}

func evalBinaryFloat(op string, left float64, rightVal Value) (Value, error) {
	var right, result float64

	switch rightVal.Type() {
	case ValueTypeFloat, ValueTypeInt:
		right = rightVal.AsFloat64()
	default:
		switch op {
		case "==", "!=":
			return &BooleanLiteral{Value: false}, nil
		default:
			return nil, errInvalidTypes
		}
	}

	switch op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return nil, errDivZero
		}
		result = left / right
	case "^":
		result = math.Pow(left, right)
	default:
		// boolean result
		return evalBinaryFloatComp(op, left, right)
	}

	return &FloatLiteral{Value: result}, nil
}

func evalBinaryFloatComp(op string, left, right float64) (Value, error) {
	var result bool

	switch op {
	case "==":
		// EQ
		result = left == right
	case "!=":
		// NE
		result = left != right
	case "<":
		// LT
		result = (left < right)
	case "<=":
		// LE
		result = (left <= right)
	case ">":
		result = (left > right)
	case ">=":
		result = (left >= right)
	default:
		return nil, errInvalidOp
	}

	return &BooleanLiteral{Value: result}, nil
}

func evalBinaryString(left, op string, rightVal Value) (Value, error) {
	switch op {
	case "+":
		// CONCAT
		result := left + rightVal.AsString()
		return &StringLiteral{Value: result}, nil
	case "==":
		result := left == rightVal.AsString()
		return &BooleanLiteral{Value: result}, nil
	case "!=":
		result := left != rightVal.AsString()
		return &BooleanLiteral{Value: result}, nil
	default:
		return nil, errInvalidOp
	}
}

func (n *UnaryExpression) String() string {
	return fmt.Sprintf("&%T{%q, %s}", *n, n.Op, n.Expr)
}

func (n *UnaryExpression) Eval() (Value, error) {
	log.Println("Unary.Eval", n)

	val, err := n.Expr.Eval()
	if err != nil {
		return nil, err
	}

	if n.Op != "!" {	
		return nil, fmt.Errorf("unknown operator: %s", n.Op)
	}

	switch val.Type() {
	case ValueTypeBool, ValueTypeInt, ValueTypeFloat:
		result := !val.AsBool()
		return &BooleanLiteral{Value: result}, nil
	default:
		return nil, errInvalidType
	}
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
