package javalanche

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	// DefaultParserTimeout praser sets default time before we reach timeout
	DefaultParserTimeout = 100 * time.Millisecond
)

// Parser represents our parser
type Parser struct {
	tokenizer *Tokenizer
	ctx       *Javalanche
	timeout   time.Duration
	tokens    []*Token
	result    ParserResult
	outCh     chan ParserResult
}

// ParserResult represents our resul struct
type ParserResult struct {
	Value Value
	Err   error
}

// NewParser represents the new parser
func NewParser(tokenizer *Tokenizer, ctx *Javalanche, timeout time.Duration) *Parser {
	if timeout < time.Millisecond {
		timeout = DefaultParserTimeout
	}

	return &Parser{
		tokenizer: tokenizer,
		ctx:       ctx,
		timeout:   timeout,
		outCh:     make(chan ParserResult),
	}
}

// Results returns the channel to watch for parser results
func (p *Parser) Results() <-chan ParserResult {
	return p.outCh
}

// IsEmpty checks if stack is empty
func (p *Parser) IsEmpty() bool {
	return len(p.tokens) == 0
}

// Pop removes a Token from the queue
func (p *Parser) Pop() *Token {
	var token *Token

	if len(p.tokens) > 0 {
		token = p.tokens[0]
		p.tokens = p.tokens[1:]
	}

	return token
}

// Push adds a non-nil Token to the stack
func (p *Parser) Push(token *Token) {

	if token != nil {
		p.tokens = append(p.tokens, token)
	}
}

// Run starts the parser
func (p *Parser) Run() {

	for {
		token, err := p.tokenizer.NextToken(p.timeout)
		switch {
		case err == os.ErrDeadlineExceeded:
			p.applyTimeout()
		case err != nil:
			done := p.applyError(err)
			if done {
				// handler wants me to end
				return
			}
		case token.Type == EOL:
			p.applyEOL()
		default:
			p.applyToken(token)
		}
	}
}

// ParseTopLevel parsest higher level of recursion
func (p *Parser) parseTopLevel() (Node, error) {
	return p.parseAssignment()
}

// ApplyError applies correct error
func (p *Parser) applyError(err error) bool {
	var terminate bool

	p.Println("applyError:", err)

	switch err {
	case io.EOF:
		// termina
		terminate = true
	}

	return terminate
}

// ApplyEOL parses tokens whe EOL token is found
func (p *Parser) applyEOL() {
	p.PrintDetails("applyEOL")

	for !p.IsEmpty() {
		token := p.Peek(0)
		var node Node
		var err error

		switch {
		case token.Type == Identifier:
			node, err = p.parseAssignment()
		case token.Type == Boolean:
			node, err = p.parseBooleanExpression()
		default:
			node, err = p.parseTopLevel()
		}

		if err != nil {
			// Fail to parse
			p.result = ParserResult{nil, err}
			return
		}

		// Call Eval directly on each Node
		value, err := node.Eval(p.ctx)

		if err != nil {
			p.result = ParserResult{nil, err}
			continue
		}

		p.result = ParserResult{value, nil}
	}
}

// ParseBooleanExpression parses bool's
func (p *Parser) parseBooleanExpression() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Operator || !(top.Value == "and" || top.Value == "&&" || top.Value == "or" || top.Value == "||") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseAnd parses logical and
func (p *Parser) parseAnd() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Keyword || !(top.Value == "and" || top.Value == "&&") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseOr parses logical or
func (p *Parser) parseOr() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Operator || !(top.Value == "or" || top.Value == "||") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseComparison parses comparison
func (p *Parser) parseComparison() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Operator || !(top.Value == "==" || top.Value == "!=" || top.Value == "<" || top.Value == ">" || top.Value == "<=" || top.Value == ">=") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseExpression parses expressions
func (p *Parser) parseExpression() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Operator || !(top.Value == "+" || top.Value == "-") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseTerm parses */
func (p *Parser) parseTerm() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Pop()
		if top == nil || top.Type != Operator || !(top.Value == "*" || top.Value == "/") {
			if top != nil {
				p.Push(top)
			}
			break
		}

		op := top.Value
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// ParseFactor parses factorail
func (p *Parser) parseFactor() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	base, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return p.parseExponent(base)
}

// ParsePrimary parses primary operations
func (p *Parser) parsePrimary() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	top := p.Pop()

	switch {
	case top.Is(Operator, "+", "-"):
		op := top.Value

		factor, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		if op == "-" {
			return &BinaryExpression{
				Op:    "*",
				Left:  NewFloat(-1),
				Right: factor,
			}, nil
		}

		return factor, nil
	case top.Is(Boolean, "true", "false"):
		value := top.Value == "true"
		return NewBoolean(value), nil
	case top.Is(Identifier):
		name := top.Value
		return &Variable{Name: name}, nil
	case top.Is(Integer):
		intVal, intErr := strconv.Atoi(top.Value)
		if intErr == nil {
			return &IntegerLiteral{Value: intVal}, nil
		}
	case top.Is(Float):
		floatVal, floatErr := strconv.ParseFloat(top.Value, 64)
		if floatErr != nil {
			return nil, floatErr
		}
		return &FloatLiteral{Value: floatVal}, nil
	case top.Is(LeftParen):
		expr, err := p.parseBooleanExpression()
		if err != nil {
			return nil, err
		}
		if top.Type != RightParen {
			return nil, errors.New("missing closing right paren")
		}
		return expr, nil
	case top.Is(Operator, "!"):
		expr, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{Op: "!", Expr: expr}, nil
	case top.Is(String):
		value := top.Value
		return &StringLiteral{Value: value}, nil
	}

	return nil, fmt.Errorf("unexpected token: %v", top.Value)
}

// ParseExponent parses exponents
func (p *Parser) parseExponent(base Node) (Node, error) {
	if p.IsEmpty() {
		return base, nil
	}
	top := p.Pop()

	for top.Is(Operator, "^") {
		op := top.Value

		exponent, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		base = &BinaryExpression{Left: base, Op: op, Right: exponent}

		if p.IsEmpty() {
			break
		}
		top = p.Pop()
	}

	if top != nil {
		p.Push(top)
	}

	return base, nil
}

// ParseAssignment parses assigments
func (p *Parser) parseAssignment() (Node, error) {
	if p.IsEmpty() {
		return nil, errors.New("unexpected end of input")
	}
	left, err := p.parseBooleanExpression()
	if err != nil {
		return nil, err
	}

	for {
		if p.IsEmpty() {
			break
		}
		top := p.Peek(1)
		if top == nil || top.Type != Operator || top.Value != "=" {
			p.Push(top)
			break
		}

		right, err := p.parseBooleanExpression()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpression{Left: left, Op: "=", Right: right}
		p.Pop()
	}

	return left, nil
}

// HasMoreTokens checks for more tokens in the tokens queue
func (p *Parser) HasMoreTokens() bool {
	return len(p.tokens) > 0
}

// ApplyTimeout appplies timeout
func (p *Parser) applyTimeout() {
	var result ParserResult

	p.PrintDetails("applyTimeout")

	switch {
	case p.IsEmpty():
		// not more tokens, emit result
		result = p.result
	default:
		result = ParserResult{nil, ErrMoreData}
	}

	p.outCh <- result
}

// Eval evaluates
func (p *Parser) Eval(_ *Javalanche) (Value, error) {
	for result := range p.outCh {
		return result.Value, result.Err
	}

	return nil, io.EOF
}

// ApplyToken pushes tokens onto quee
func (p *Parser) applyToken(token *Token) {
	p.Push(token)
}

// Peek peeks into next token on the queue
func (p *Parser) Peek(index int) *Token {
	var token *Token
	var i int

	switch {
	case index < 0:
		// reverse
		i = len(p.tokens) + index
		if i >= 0 {
			token = p.tokens[i]
		}
	case index < len(p.tokens):
		// found
		i = index
		token = p.tokens[index]
	}

	return token
}
