package javalanche

import (
	"fmt"
)

var (
	_ Value          = (*BooleanLiteral)(nil)
	_ Node           = (*BooleanLiteral)(nil)
	_ fmt.GoStringer = (*BooleanLiteral)(nil)
	_ fmt.Stringer   = (*BooleanLiteral)(nil)
)

type BooleanLiteral struct {
	Value bool
}

func NewBoolean(v bool) *BooleanLiteral {
	return &BooleanLiteral{Value: v}
}

func (n *BooleanLiteral) GoString() string {
	return fmt.Sprintf("NewBoolean(%s)", n.String())
}

func (n *BooleanLiteral) String() string {
	if n.Value {
		return "true"
	}
	return "false"
}

func (n *BooleanLiteral) Type() ValueType {
	return ValueTypeBool
}

func (n *BooleanLiteral) AsFloat64() float64 {
	if n.Value {
		return 1
	}
	return 0
}

func (n *BooleanLiteral) AsBool() bool {
	return n.Value
}

func (n *BooleanLiteral) AsString() string {
	return n.String()
}

func (n *BooleanLiteral) Eval() (Value, error) {
	return n, nil
}

// Equal attempts to apply the == operation to
// this boolean and a given right-value
func (n *BooleanLiteral) Equal(v Value) bool {
	if m, ok := v.(*BooleanLiteral); ok {
		return n.Value == m.Value
	}
	return false
}
