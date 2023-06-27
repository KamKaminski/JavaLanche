package javalanche

import (
	"io"
	"os"
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
	stage     Stage
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

// IsEmpty checks if the stage is empty
func (p *Parser) IsEmpty() bool {
	return p.stage.IsEmpty()
}

// GetPrecedence function to get the precedence of a token
func getOperatorPrecedence(op string) int {
	switch op {
	case "=":
		return 1
	case "or", "||":
		return 2
	case "and", "&&":
		return 3
	case "==":
		return 4
	case "!=":
		return 5
	case "<", ">", "<=", ">=":
		return 6
	case "+":
		return 7
	case "-":
		return 8
	case "*", "/", "%":
		return 9
	case "^":
		return 10
	case "!":
		return 11
	default:
		return 12
	}
}

func findHighestPrecedenceOperatorInRange(nodes []StageNode, start, end int) (int, *Token) {
	pivot, op := findHighestPrecedenceOperator(nodes[start:end])
	if op != nil {
		return pivot + start, op
	}
	return -1, nil
}

func findHighestPrecedenceOperator(nodes []StageNode) (int, *Token) {
	var op *Token
	var maxPrecedence = -1
	var maxPrecedenceIndex = -1

	for i, node := range nodes {
		if token, ok := node.Token(); ok {
			if token.Type == Operator {
				precedence := getOperatorPrecedence(token.Value)
				if precedence >= maxPrecedence {
					// last wins, because some binary operators
					// are also prefixed unary operators
					maxPrecedence = precedence
					maxPrecedenceIndex = i
					op = token
				}
			}
		}
	}

	return maxPrecedenceIndex, op
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

// ApplyError applies correct error
func (p *Parser) applyError(err error) bool {
	var terminate bool

	p.Println("applyError:", err)

	switch err {
	case io.EOF:
		// terminate
		terminate = true
	}

	return terminate
}

// ApplyEOL parses tokens whe EOL token is found
func (p *Parser) applyEOL() {
	p.PrintDetails("applyEOL")

	node, err := p.stage.Parse()
	if err != nil {
		// Fail to parse
		p.Println("applyEOL:", "Stage.Parse:", "err:", err)
		p.result = ParserResult{nil, err}
		return
	}

	p.Println("applyEOL:", "Stage.Parse:", node)

	// Call Eval directly on each Node
	value, err := node.Eval(p.ctx)
	switch {
	case err != nil:
		p.Println("applyEOL:", node, "→ err:", err)
		p.result = ParserResult{nil, err}
	default:
		p.Println("applyEOL:", node, "→", value)
		p.result = ParserResult{value, nil}
	}
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
	if token != nil {
		p.stage.AppendTokens(token)
	}
}
