package javalanche

import (
	"errors"
)

var (
	errInvalidTypes = errors.New("invalid types")
	errDivZero      = errors.New("division by zero")
)

// AddValuer provides add interface
type AddValuer interface {
	AddValue(Value) (Value, error)
}

// AndValuer provides And interface
type AndValuer interface {
	AndValue(Value) (Value, error)
}

// SubValuer provides Sub interface
type SubValuer interface {
	SubValue(Value) (Value, error)
}

// OrValuer provides or interface
type OrValuer interface {
	OrValue(Value) (Value, error)
}

// MulValuer provides Mul interface
type MulValuer interface {
	MulValue(Value) (Value, error)
}

// DivValuer provides div interface
type DivValuer interface {
	DivValue(Value) (Value, error)
}

type NegValuer interface {
	NegValue() (Value, error)
}

// UpValuer provides up interface
type UpValuer interface {
	UpValue(Value) (Value, error)
}

// GreaterValuer provides > interface
type GreaterValuer interface {
	GreaterValue(Value) (Value, error)
}

// LesserValuer provides < interface
type LesserValuer interface {
	LesserValue(Value) (Value, error)
}

// GreaterEqualValuer provides >= interface
type GreaterEqualValuer interface {
	GreaterEqualValue(Value) (Value, error)
}

// LesserEqualValuer provides <= interface
type LesserEqualValuer interface {
	LesserEqualValue(Value) (Value, error)
}

// LogicalNotValuer represents a Value that
// implements Unary negation.
type LogicalNotValuer interface {
	LogicalNotValue() (Value, error)
}

// Value represents value interface
type Value interface {
	Type() ValueType
	AsFloat64() float64
	AsString() string
	AsBool() bool

	// Ops
	Equal(Value) bool
}

// Node represents nodes interface
type Node interface {
	Eval(ctx *Javalanche) (Value, error)
}

// ValueType represents the type of a Value
type ValueType int

// ValueTypeUnknown indicates the Type wasn't set
// ValueTypeInt indicates the Value contains an Integer
// ValueTypeFloat indicates the Value contains an Float
// ValueTypeString inidcates the Value contains String
// ValueTypeBool indicates the Value contains Bool
const (
	ValueTypeUnknown ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeString
	ValueTypeBool
)
