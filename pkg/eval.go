// Package javalanche provides an interpreter of the javalanche language
package javalanche

import (
	"fmt"
	"strings"
)

// Evaluator struct represents map of the variables
type Evaluator struct {
	Variable map[string]Value
}

// NewEvaluator creates a new Evaluator with an empty variable map
func NewEvaluator() *Evaluator {
	return &Evaluator{
		Variable: make(map[string]Value),
	}
}

// SetValue Assigns value to given variable
func (e *Evaluator) SetValue(name string, v Value) error {
	e.Variable[name] = v
	return nil
}

// GetValue retrieves Value of given variable
func (e *Evaluator) GetValue(name string) (Value, error) {
	if v, ok := e.Variable[name]; ok {
		return v, nil

	}
	return nil, fmt.Errorf("variable %q not found", name)
}

// EvalString evaluates expressions
func EvalString(evaluator *Evaluator, exprs ...string) (Value, error) {
	var res Value

	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if expr != "" {
			var err error

			res, err = evalSingleExpression(evaluator, expr)
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}

func evalSingleExpression(evaluator *Evaluator, expr string) (Value, error) {
	// Since Tokenizer implements the idea of I/O, we need to turn our input into the reader
	reader := strings.NewReader(expr)
	lexer := NewTokenizer(reader)
	parser := NewParser(lexer, evaluator)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	value, err := ast.Eval(evaluator)
	if err != nil {
		return nil, err
	}

	return value, nil
}
