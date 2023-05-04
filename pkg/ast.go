package javalanche

import (
	"errors"
)

var (
	errInvalidTypes = errors.New("invalid types")
	//	errInvalidOp    = errors.New("invalid operator")
	errDivZero = errors.New("division by zero")
)

type AddValuer interface {
	AddValue(Value) (Value, error)
}

type AndValuer interface {
	AndValue(Value) (Value, error)
}

type SubValuer interface {
	SubValue(Value) (Value, error)
}

type OrValuer interface {
	OrValue(Value) (Value, error)
}

type MulValuer interface {
	MulValue(Value) (Value, error)
}

type DivValuer interface {
	DivValue(Value) (Value, error)
}

type UpValuer interface {
	UpValue(Value) (Value, error)
}

type GreaterValuer interface {
	GreaterValue(Value) (Value, error)
}

type LesserValuer interface {
	LesserValue(Value) (Value, error)
}

type GreaterEqualValuer interface {
	GreaterEqualValue(Value) (Value, error)
}

type LesserEqualValuer interface {
	LesserEqualValue(Value) (Value, error)
}

// LogicalNotValue represents a Value that
// implements Unary negation.
type LogicalNotValuer interface {
	LogicalNotValue() (Value, error)
}

type Value interface {
	Type() ValueType
	AsFloat64() float64
	AsString() string
	AsBool() bool

	// Ops
	Equal(Value) bool
}

type Node interface {
	Eval(ctx *Evaluator) (Value, error)
}

type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeString
	ValueTypeBool
)
