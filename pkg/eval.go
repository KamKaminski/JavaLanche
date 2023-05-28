// Package javalanche provides an interpreter of the javalanche language
package javalanche

import (
	"fmt"
)

// SetValue Assigns value to given variable
func (ctx *Javalanche) SetValue(name string, v Value) error {
	ctx.Variable[name] = v
	return nil
}

// GetValue retrieves Value of given variable
func (ctx *Javalanche) GetValue(name string) (Value, error) {
	if v, ok := ctx.Variable[name]; ok {
		return v, nil

	}
	return nil, fmt.Errorf("variable %q not found", name)
}

// EvalLine evaluates expression lines
func (ctx *Javalanche) EvalLine(lines ...string) (Value, error) {
	if err := ctx.ParseLine(lines...); err != nil {
		return nil, err
	}

	return ctx.Eval()
}

// Eval returns the latest result or semantic error
func (ctx *Javalanche) Eval() (Value, error) {
	return ctx.parser.Eval(ctx)
}
