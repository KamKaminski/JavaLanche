// Package javalanche provides an interpreter of the javalanche language
package javalanche

import (
	"errors"
)

var (
	// ErrMoreData provides new error when more data is needed to evalute
	ErrMoreData = errors.New("more data needed to evaluate the statement")
)

// Javalanche represts Interpreter for Javalanche language
type Javalanche struct {
	Evaluator *Evaluator
}

// New Creates the new instance of Javalanche
func New() *Javalanche {
	return &Javalanche{
		Evaluator: NewEvaluator(),
	}
}

// EvalLine evaluates a line of input
func (j *Javalanche) EvalLine(s string) (Value, error) {
	return EvalString(j.Evaluator, s)
}
