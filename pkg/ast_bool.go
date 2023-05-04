package javalanche

import (
	"fmt"
)

var (
	_ Value            = (*BooleanLiteral)(nil)
	_ Node             = (*BooleanLiteral)(nil)
	_ fmt.GoStringer   = (*BooleanLiteral)(nil)
	_ fmt.Stringer     = (*BooleanLiteral)(nil)
	_ UpValuer         = (*BooleanLiteral)(nil)
	_ AndValuer        = (*BooleanLiteral)(nil)
	_ OrValuer         = (*BooleanLiteral)(nil)
	_ LogicalNotValuer = (*BooleanLiteral)(nil)
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

func (n *BooleanLiteral) Eval(ctx *Evaluator) (Value, error) {
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

func (n *BooleanLiteral) UpValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *BooleanLiteral:
		// n XOR v
		var res bool

		if (n.Value || right.Value) &&
			!(n.Value && right.Value) {
			res = true
		}

		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *BooleanLiteral) OrValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *BooleanLiteral:
		res := n.Value || right.Value
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *BooleanLiteral) AndValue(v Value) (Value, error) {
	switch right := v.(type) {
	case *BooleanLiteral:
		res := n.Value && right.Value
		return NewBoolean(res), nil
	default:
		return nil, errInvalidTypes
	}
}

func (n *BooleanLiteral) LogicalNotValue() (Value, error) {
	return NewBoolean(!n.Value), nil
}
