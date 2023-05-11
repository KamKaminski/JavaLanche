// Package javalanche provides an interpreter of the javalanche language
package javalanche

import (
	"fmt"
	"strings"
)

// SetValue Assigns value to given variable
func (e *Javalanche) SetValue(name string, v Value) error {
	e.Variable[name] = v
	return nil
}

// GetValue retrieves Value of given variable
func (e *Javalanche) GetValue(name string) (Value, error) {
	if v, ok := e.Variable[name]; ok {
		return v, nil

	}
	return nil, fmt.Errorf("variable %q not found", name)
}

// EvalString evaluates expressions
func (ctx *Javalanche) EvalLine(exprs ...string) (Value, error) {
	var res Value

	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if expr != "" {
			var err error

			res, err = evalSingleExpression(ctx, expr)
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil

}

func evalSingleExpression(ctx *Javalanche, expr string) (Value, error) {
	// Since Tokenizer implements the idea of I/O, we need to turn our input into the reader
	reader := strings.NewReader(expr)
	lexer := NewTokenizer(reader)
	parser := NewParser(lexer, ctx)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	value, err := ast.Eval(ctx)
	if err != nil {
		return nil, err
	}

	return value, nil
}
