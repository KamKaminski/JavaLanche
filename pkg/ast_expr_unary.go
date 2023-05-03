package javalanche

import (
	"fmt"
	"log"
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
	switch n.Op {
	case "!":
		// prefix
		return fmt.Sprintf("%s%s", n.Op, n.Expr)
	default:
		// suffix
		return fmt.Sprintf("%s%s", n.Expr, n.Op)
	}
}

func (n *UnaryExpression) Eval() (Value, error) {
	log.Printf("Eval: %#v", n)

	val, err := n.Expr.Eval()
	if err != nil {
		return nil, err
	}

	switch n.Op {
	case "!":
		if left, ok := val.(LogicalNotValuer); ok {
			return left.LogicalNotValue()
		}
	}

	err = fmt.Errorf("operator %q can't be used on %s", n.Op, val)
	return nil, err
}
