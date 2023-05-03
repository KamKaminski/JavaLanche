package javalanche

import (
	"fmt"
	"log"
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

	switch n.Op {
	case "==":
		eq := leftVal.Equal(rightVal)
		return NewBoolean(eq), nil
	case "!=":
		eq := leftVal.Equal(rightVal)
		return NewBoolean(!eq), nil
	case "&&", "and":
		// left.AndValue(right)
		// left is of type AndValuer,
		//we can only use it if the cast succeeded (ok)
		if left, ok := leftVal.(AndValuer); ok {
			return left.AndValue(rightVal)
		}
	case "||", "or":
		// left.OrValue(right)
		if left, ok := leftVal.(OrValuer); ok {
			return left.OrValue(rightVal)
		}
	case "^":
		// left.UpValue(right)
		if left, ok := leftVal.(UpValuer); ok {
			return left.UpValue(rightVal)
		}
	case "+":
		if left, ok := leftVal.(AddValuer); ok {
			return left.AddValue(rightVal)
		}
	case "-":
		if left, ok := leftVal.(SubValuer); ok {
			return left.SubValue(rightVal)
		}
	case "/":
		if left, ok := leftVal.(DivValuer); ok {
			return left.DivValue(rightVal)
		}
	case "*":
		if left, ok := leftVal.(MulValuer); ok {
			return left.MulValue(rightVal)
		}
	case ">":
		if left, ok := leftVal.(GreaterValuer); ok {
			return left.GreaterValue(rightVal)
		}
	case "<":
		if left, ok := leftVal.(LesserValuer); ok {
			return left.LesserValue(rightVal)
		}
	case ">=":
		if left, ok := leftVal.(GreaterEqualValuer); ok {
			return left.GreaterEqualValue(rightVal)
		}
	case "<=":
		if left, ok := leftVal.(LesserEqualValuer); ok {
			return left.LesserEqualValue(rightVal)
		}
	}

	err = fmt.Errorf("operator %q can't be used on %s", n.Op, leftVal)
	return nil, err
}
