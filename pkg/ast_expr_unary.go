package javalanche

import (
	"fmt"
)

var (
	_ Node           = (*UnaryExpression)(nil)
	_ fmt.GoStringer = (*UnaryExpression)(nil)
	_ fmt.Stringer   = (*UnaryExpression)(nil)
)

type UnaryExpression struct {
	Op   string
	Expr Node
}

func (n *UnaryExpression) GoString() string {
	return fmt.Sprintf("&UnaryExpression{%q, %#v}", n.Op, n.Expr)
}

func (n *UnaryExpression) String() string {
	switch {
	case isPrefixUnaryOperator(n.Op):
		// prefix
		return fmt.Sprintf("%s%s", n.Op, n.Expr)
	default:
		// suffix
		return fmt.Sprintf("%s%s", n.Expr, n.Op)
	}
}

func (n *UnaryExpression) Eval(ctx *Javalanche) (Value, error) {
	val, err := n.Expr.Eval(ctx)
	if err != nil {
		// bad operand
		return nil, err
	}

	switch n.Op {
	case "++":
		if v, ok := val.(AddValuer); ok {
			val, err = v.AddValue(NewInteger(1))
			if err != nil {
				return nil, err
			}
			return nil, setValue(ctx, n.Expr, val)
		}
	case "--":
		if v, ok := val.(SubValuer); ok {
			val, err = v.SubValue(NewInteger(1))
			if err != nil {
				return nil, err
			}
			return nil, setValue(ctx, n.Expr, val)
		}
	case "!":
		if left, ok := val.(LogicalNotValuer); ok {
			return left.LogicalNotValue()
		}
	}

	err = fmt.Errorf("operator %q can't be used on %s", n.Op, val)
	return nil, err
}

func setValue(ctx *Javalanche, node Node, val Value) error {
	if operand, ok := node.(SetValuer); ok {
		return operand.SetValue(ctx, val)
	}

	return fmt.Errorf("%s can't be set to %s", node, val)
}
