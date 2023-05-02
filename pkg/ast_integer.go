package javalanche

import (
	"fmt"
	"log"
)

var (
	_ Value          = (*IntegerLiteral)(nil)
	_ Node           = (*IntegerLiteral)(nil)
	_ fmt.GoStringer = (*IntegerLiteral)(nil)
	_ fmt.Stringer   = (*IntegerLiteral)(nil)
)

type IntegerLiteral struct {
	Value int
}

func NewInteger(n int) *IntegerLiteral {
	return &IntegerLiteral{Value: n}
}

func (n *IntegerLiteral) GoString() string {
	return fmt.Sprintf("NewInteger(%v)", n.Value)
}

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n *IntegerLiteral) Type() ValueType {
	return ValueTypeInt
}

func (n *IntegerLiteral) AsFloat64() float64 {
	return float64(n.Value)
}

func (n *IntegerLiteral) AsString() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n *IntegerLiteral) AsBool() bool {
	return n.Value != 0
}

func (n *IntegerLiteral) Eval() (Value, error) {
	log.Printf("Eval: %#v", n)
	return n, nil
}

func (n *IntegerLiteral) Equal(v Value) bool {
	if m, ok := v.(*IntegerLiteral); ok {
		return n.Value == m.Value
	}
	if m, ok := v.(*FloatLiteral); ok {
		return float64(n.Value) == m.Value
	}

	return false
}
