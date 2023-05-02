package javalanche

import (
	"errors"
)

var (
	errInvalidType  = errors.New("invalid type")
	errInvalidTypes = errors.New("invalid types")
	errInvalidOp    = errors.New("invalid operator")
	errDivZero      = errors.New("division by zero")
)

type Value interface {
	Type() ValueType
	AsFloat64() float64
	AsString() string
	AsBool() bool

	// Ops
	Equal(Value) bool
}

type Node interface {
	Eval() (Value, error)
}

type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeString
	ValueTypeBool
)
