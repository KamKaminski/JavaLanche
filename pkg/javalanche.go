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
	Variable map[string]Value
}

// New Creates the new instance of Javalanche
func New() *Javalanche {
	return &Javalanche{
		Variable: make(map[string]Value),
	}
}
