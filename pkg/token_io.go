package javalanche

import (
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

// Err log
var (
	errLexEmitNoArgs = errors.New("emit called without argument")
)

// Keywords list
var keywords = []string{
	"if",
	"else",
	"for",
	"while",
	"funx",
	"return",
	"var",
	"true",
	"false",
	"int",
	"float",
	"bool",
	"string",
	"const",
	"print",
	"in",
	"let",
}

// Tokenizer represents tokenizer
type Tokenizer struct {
	reader *Reader
	buffer []Token

	outCh  chan *TokenResult
	closed bool
}

// Close closes channel
func (t *Tokenizer) Close() error {
	if !t.closed {
		close(t.outCh)
		t.closed = true
	}
	return nil
}

// Results Channels result to target channel
func (t *Tokenizer) Results() <-chan *TokenResult {
	return t.outCh
}

// NextToken reads the output channel to return the next token
// or error
func (t *Tokenizer) NextToken(deadline time.Duration) (*Token, error) {
	select {
	case result, ok := <-t.outCh:
		switch {
		case !ok:
			return nil, io.EOF
		case result.Token != nil:
			return result.Token, nil
		case result.Err != nil:
			return nil, result.Err
		default:
			// nil, nil
			panic("unreachable")
		}
	case <-time.After(deadline):
		return nil, os.ErrDeadlineExceeded
	}
}

// emit emits tokens
func (t *Tokenizer) emit(res *TokenResult) {
	switch {
	case res == nil:
		panic(errLexEmitNoArgs)
	case t.closed:
		// ignore
	default:
		t.outCh <- res
	}
}

// emitErrors emits errors
func (t *Tokenizer) emitError(err error) {
	switch {
	case err == nil:
		panic(errLexEmitNoArgs)
	default:
		res := &TokenResult{
			Err: err,
		}

		t.emit(res)
	}
}

// emitValue emits value of the toen and type
func (t *Tokenizer) emitValue(typ TokenType, val string) {
	res := &TokenResult{
		Token: &Token{
			Type:  typ,
			Value: val,
		},
	}

	t.emit(res)
}

// emitToken emits whole token
func (t *Tokenizer) emitToken(typ TokenType) {
	t.emitValue(typ, t.reader.EmitString())
}

// NewTokenizer represents newtokenizer
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: NewReader(r),
		outCh:  make(chan *TokenResult),
	}
}

// PushBack gives us ability to pushback
func (t *Tokenizer) PushBack(token *Token) {
	t.buffer = append(t.buffer, *token)
}

// Run starts tokenizer
func (t *Tokenizer) Run() {
	defer t.Close()

	for fn := lexText; fn != nil; {
		fn = fn(t)
	}
}

// accept checks whether rune is valid
func (t *Tokenizer) accept(valid string) bool {
	return t.reader.Accept(func(r rune) bool {
		return strings.ContainsRune(valid, r)
	})
}

// acceptFn uses helper functions to check whether rune is valid
func (t *Tokenizer) acceptFn(match func(rune) bool) bool {
	return t.reader.Accept(match)
}

// acceptAll validates run of runes
func (t *Tokenizer) acceptAll(valid string) bool {
	return t.reader.AcceptAll(func(r rune) bool {
		return strings.ContainsRune(valid, r)
	})
}

// acceptAllFn nvalidates run of runes
func (t *Tokenizer) acceptAllFn(match func(rune) bool) bool {
	return t.reader.AcceptAll(match)
}
