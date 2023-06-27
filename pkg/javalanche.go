// Package javalanche provides an interpreter of the javalanche language
package javalanche

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	// ErrMoreData provides new error when more data is needed to evalute
	ErrMoreData = errors.New("more data needed to evaluate the statement")
)

// Javalanche represts Interpreter for Javalanche language
type Javalanche struct {
	Variable map[string]Value

	buf    bytes.Buffer
	mu     sync.Mutex
	lexer  *Tokenizer
	parser *Parser
}

// New Creates the new instance of Javalanche
func New() *Javalanche {
	ctx := &Javalanche{
		Variable: make(map[string]Value),
	}

	timeout := 100 * time.Millisecond
	ctx.lexer = NewTokenizer(&ctx.buf)
	ctx.parser = NewParser(ctx.lexer, ctx, timeout)

	go ctx.lexer.Run()
	go ctx.parser.Run()

	return ctx
}

// ParseLine feeds the parser with a new line of javalanche
func (ctx *Javalanche) ParseLine(lines ...string) error {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			ctx.buf.WriteString(line)
			ctx.buf.WriteRune('\n')
		}
	}

	return nil
}
