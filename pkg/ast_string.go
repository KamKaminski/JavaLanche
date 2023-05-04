package javalanche

import (
	"fmt"
	"log"
	"strconv"
)

var (
	_ Value          = (*StringLiteral)(nil)
	_ Node           = (*StringLiteral)(nil)
	_ fmt.GoStringer = (*StringLiteral)(nil)
	_ fmt.Stringer   = (*StringLiteral)(nil)
	_ AddValuer      = (*StringLiteral)(nil)
)

type StringLiteral struct {
	Value string
}

func NewString(s string) *StringLiteral {
	return &StringLiteral{Value: s}
}

func (n *StringLiteral) GoString() string {
	return fmt.Sprintf("NewString(%q)", n.Value)
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

func (n *StringLiteral) Eval(ctx *Evaluator) (Value, error) {
	log.Println("String.Eval", n.Value)

	return n, nil
}

func (n *StringLiteral) Equal(v Value) bool {
	return n.Value == v.AsString()
}

// AddValue concatenates two strings, converting the argument
// to string first if needed.
func (n *StringLiteral) AddValue(v Value) (Value, error) {
	rightStr := v.AsString()
	result := n.Value + rightStr
	return NewString(result), nil
}
