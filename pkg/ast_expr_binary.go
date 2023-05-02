package javalanche

import (
	"fmt"
	"log"
	"math"
)

var (
	_ Node           = (*BinaryExpression)(nil)
	_ fmt.GoStringer = (*BinaryExpression)(nil)
	_ fmt.Stringer   = (*BinaryExpression)(nil)
)

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}

func (n *BinaryExpression) GoString() string {
	return fmt.Sprintf("&BinaryExpression{%#v, %q, %#v}", n.Left, n.Op, n.Right)
}

func (n *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", n.Left, n.Op, n.Right)
}

func (n *BinaryExpression) Eval() (Value, error) {
	log.Printf("Eval: %#v", n)

	leftVal, err := n.Left.Eval()
	if err != nil {
		return nil, err
	}

	rightVal, err := n.Right.Eval()
	if err != nil {
		return nil, err
	}

	switch leftVal.Type() {
	case ValueTypeBool:
		return evalBinaryBool(n.Op, leftVal.AsBool(), rightVal)
	case ValueTypeString:
		return evalBinaryString(n.Op, leftVal.AsString(), rightVal)
	case ValueTypeFloat, ValueTypeInt:
		return evalBinaryFloat(n.Op, leftVal.AsFloat64(), rightVal)
	default:
		return evalBinaryFloat(n.Op, leftVal.AsFloat64(), rightVal)
	}
}

func evalBinaryBool(op string, left bool, rightVal Value) (Value, error) {
	var right, result bool

	switch {
	case rightVal.Type() == ValueTypeBool:
		right = rightVal.AsBool()
	case op == "==", op == "!=":
		return &BooleanLiteral{Value: false}, nil
	default:
		return nil, errInvalidTypes
	}

	switch op {
	case "&&", "and":
		// AND
		result = left && right
	case "||", "or":
		// OR
		result = left || right
	case "==":
		// EQ
		result = left == right
	case "!=":
		// NE
		result = left != right
	case "^":
		// XOR
		result = (left && !right) || (!left && right)
	default:
		return nil, errInvalidOp
	}

	return &BooleanLiteral{Value: result}, nil
}

func evalBinaryFloat(op string, left float64, rightVal Value) (Value, error) {
	var right, result float64

	switch rightVal.Type() {
	case ValueTypeFloat, ValueTypeInt:
		right = rightVal.AsFloat64()
	default:
		switch op {
		case "==", "!=":
			return &BooleanLiteral{Value: false}, nil
		default:
			return nil, errInvalidTypes
		}
	}

	switch op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return nil, errDivZero
		}
		result = left / right
	case "^":
		result = math.Pow(left, right)
	default:
		// boolean result
		return evalBinaryFloatComp(op, left, right)
	}

	return &FloatLiteral{Value: result}, nil
}

func evalBinaryFloatComp(op string, left, right float64) (Value, error) {
	var result bool

	switch op {
	case "==":
		// EQ
		result = left == right
	case "!=":
		// NE
		result = left != right
	case "<":
		// LT
		result = (left < right)
	case "<=":
		// LE
		result = (left <= right)
	case ">":
		result = (left > right)
	case ">=":
		result = (left >= right)
	default:
		return nil, errInvalidOp
	}

	return &BooleanLiteral{Value: result}, nil
}

func evalBinaryString(op, left string, rightVal Value) (Value, error) {
	switch op {
	case "+":
		// CONCAT
		result := left + rightVal.AsString()
		return NewString(result), nil
	case "==":
		result := left == rightVal.AsString()
		return &BooleanLiteral{Value: result}, nil
	case "!=":
		result := left != rightVal.AsString()
		return &BooleanLiteral{Value: result}, nil
	default:
		return nil, errInvalidOp
	}
}
