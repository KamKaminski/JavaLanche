package javalanche

import (
	"fmt"
	"math"
	"strconv"
)

var (
	_ Value              = (*FloatLiteral)(nil)
	_ Node               = (*FloatLiteral)(nil)
	_ fmt.GoStringer     = (*FloatLiteral)(nil)
	_ fmt.Stringer       = (*FloatLiteral)(nil)
	_ AddValuer          = (*FloatLiteral)(nil)
	_ SubValuer          = (*FloatLiteral)(nil)
	_ MulValuer          = (*FloatLiteral)(nil)
	_ DivValuer          = (*FloatLiteral)(nil)
	_ UpValuer           = (*FloatLiteral)(nil)
	_ GreaterEqualValuer = (*FloatLiteral)(nil)
	_ LesserEqualValuer  = (*FloatLiteral)(nil)
	_ LesserValuer       = (*FloatLiteral)(nil)
	_ GreaterValuer      = (*FloatLiteral)(nil)
)

type FloatLiteral struct {
	Value float64
}

func NewFloat(n float64) *FloatLiteral {
	return &FloatLiteral{Value: n}
}

// NewFloatString converts string value into float value
func NewFloatString(s string) (*FloatLiteral, error) {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return &FloatLiteral{Value: floatVal}, nil
	}

	return nil, &ErrInvalidValue{s}
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

func (n *FloatLiteral) Eval(ctx *Javalanche) (Value, error) {
	return n, nil
}

func (n *FloatLiteral) Equal(v Value) bool {
	switch right := v.(type) {
	case *FloatLiteral:
		return n.Value == right.Value
	case *IntegerLiteral:
		return n.Value == float64(right.Value)
	default:
		return false
	}
}

func (n *FloatLiteral) AddValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := n.Value + right.Value
		return NewFloat(res), nil
	case *IntegerLiteral:
		res := n.Value + float64(right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) SubValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := (n.Value) - (right.Value)
		return NewFloat(res), nil
	case *IntegerLiteral:
		res := (n.Value) - (float64)(right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) DivValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		if right.Value == 0 {
			return nil, errDivZero
		}
		// float/float -> float
		res := n.Value / right.Value
		return NewFloat(res), nil
	case *IntegerLiteral:
		if right.Value == 0 {
			return nil, errDivZero
		}
		// float / int -> float
		res := n.Value / (float64)(right.Value)
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) MulValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		// 3.5 * 2.5 -> 8.75
		res := n.Value * right.Value
		return NewFloat(res), nil
	case *IntegerLiteral:
		res := n.Value * (float64)(right.Value)
		// 3.5 * 1.0 -> 3.5
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) UpValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := math.Pow(n.Value, right.Value)
		return NewFloat(res), nil
	case *IntegerLiteral:
		res := math.Pow(n.Value, (float64)(right.Value))
		return NewFloat(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) LesserValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := (n.Value < right.Value)
		return NewBoolean(res), nil
	case *IntegerLiteral:
		res := (n.Value < (float64)(right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) GreaterValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := (n.Value > right.Value)
		return NewBoolean(res), nil
	case *IntegerLiteral:
		res := (n.Value > (float64)(right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) GreaterEqualValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := (n.Value >= right.Value)
		return NewBoolean(res), nil
	case *IntegerLiteral:
		res := (n.Value >= (float64)(right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *FloatLiteral) LesserEqualValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *FloatLiteral:
		res := (n.Value <= right.Value)
		return NewBoolean(res), nil
	case *IntegerLiteral:
		res := (n.Value <= (float64)(right.Value))
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}
