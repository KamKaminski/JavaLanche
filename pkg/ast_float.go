package javalanche

import (
	"fmt"
)

var (
	_ Value          = (*FloatLiteral)(nil)
	_ Node           = (*FloatLiteral)(nil)
	_ fmt.GoStringer = (*FloatLiteral)(nil)
	_ fmt.Stringer   = (*FloatLiteral)(nil)
)

type FloatLiteral struct {
	Value float64
}

func NewFloat(n float64) *FloatLiteral {
	return &FloatLiteral{Value: n}
}

func (n *FloatLiteral) GoString() string {
	return fmt.Sprintf("NewFloat(%f)", n.Value)
}

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *FloatLiteral) Type() ValueType {
	return ValueTypeFloat
}

func (n *FloatLiteral) AsFloat64() float64 {
	return n.Value
}

func (n *FloatLiteral) AsString() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *FloatLiteral) AsBool() bool {
	return n.Value != 0
}

func (n *FloatLiteral) Eval() (Value, error) {
	return n, nil
}

func (n *FloatLiteral) Equal(v Value) bool {
	return false
}
