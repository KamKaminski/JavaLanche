package javalanche

import (
	"fmt"
	"log"
	"math"
)

var (
	_ Value             = (*IntegerLiteral)(nil)
	_ Node              = (*IntegerLiteral)(nil)
	_ fmt.GoStringer    = (*IntegerLiteral)(nil)
	_ fmt.Stringer      = (*IntegerLiteral)(nil)
	_ AddValuer         = (*IntegerLiteral)(nil)
	_ SubValuer         = (*IntegerLiteral)(nil)
	_ MulValuer         = (*IntegerLiteral)(nil)
	_ DivValuer         = (*IntegerLiteral)(nil)
	_ UpValuer          = (*IntegerLiteral)(nil)
	_ GreaterValuer     = (*IntegerLiteral)(nil)
	_ LesserEqualValuer = (*IntegerLiteral)(nil)
	_ LesserValuer      = (*IntegerLiteral)(nil)
	_ GreaterValuer     = (*IntegerLiteral)(nil)
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
func (n *IntegerLiteral) AddValue(v Value) (Value, error) {
	// if v is of the specified type
	// the switch will test the supported types and use the best match
	switch right := v.(type) {
	case *IntegerLiteral:
		res := n.Value + right.Value
		return NewInteger(res), nil
	case *FloatLiteral:
		res := (float64)(n.Value) + right.Value
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) SubValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := n.Value - right.Value
		return NewInteger(res), nil
	case *FloatLiteral:
		res := (float64)(n.Value) - (right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) DivValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		if right.Value == 0 {
			return nil, errDivZero
		}
		res := (float64)(n.Value) / (float64)(right.Value)
		return NewFloat(res), nil
	case *FloatLiteral:
		if right.Value == 0 {
			return nil, errDivZero
		}
		res := (float64)(n.Value) / (right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) MulValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := n.Value * right.Value
		return NewInteger(res), nil
	case *FloatLiteral:
		res := (float64)(n.Value) * (right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) UpValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		// Using primivite as math.pow signature requires floats
		res := (n.Value ^ right.Value)
		return NewInteger(res), nil
	case *FloatLiteral:
		// floats values so we can use math.pow
		res := math.Pow((float64)(n.Value), right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) LesserValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := (n.Value < right.Value)
		return NewBoolean(res), nil
	case *FloatLiteral:
		res := ((float64)(n.Value) < (right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) GreaterValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := (n.Value > right.Value)
		return NewBoolean(res), nil
	case *FloatLiteral:
		res := ((float64)(n.Value) > (right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) GreaterEqualValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := (n.Value >= right.Value)
		return NewBoolean(res), nil
	case *FloatLiteral:
		res := ((float64)(n.Value) >= (right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *IntegerLiteral) LesserEqualValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *IntegerLiteral:
		res := (n.Value <= right.Value)
		return NewBoolean(res), nil
	case *FloatLiteral:
		res := ((float64)(n.Value) <= (right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}
