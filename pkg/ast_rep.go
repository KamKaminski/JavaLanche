package javalanche

var (
	_ Node = (*BodyNode)(nil)
	_ Node = (*IfElseNode)(nil)
	_ Node = (*ForNode)(nil)
)

// BodyNode represents our body node
type BodyNode []Node

// Eval evaluates body node
func (body BodyNode) Eval(ctx *Javalanche) (Value, error) {
	var val Value
	var err error

	for _, n := range body {
		val, err = n.Eval(ctx)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

// IfElseNode is struct of IfElse Node
type IfElseNode struct {
	Condition Node
	TrueBody  Node
	FalseBody Node
}

// Eval evaluates if/elif/else
func (n *IfElseNode) Eval(ctx *Javalanche) (Value, error) {
	condVal, err := n.Condition.Eval(ctx)
	switch {
	case err != nil:
		return nil, err
	case condVal.AsBool():
		return n.TrueBody.Eval(ctx)
	case n.FalseBody != nil:
		return n.FalseBody.Eval(ctx)
	default:
		return nil, nil
	}
}

// ForNode represents for loops strucct
type ForNode struct {
	Condition Node
	Body      Node
}

// Eval evaluates for loop
func (n *ForNode) Eval(ctx *Javalanche) (Value, error) {
	var val Value

	for {
		condVal, err := n.Condition.Eval(ctx)
		switch {
		case err != nil:
			// failed to evaluate condition
			return nil, err
		case !condVal.AsBool():
			// break
			return val, nil
		case n.Body == nil:
			// no body, break to prevent infinite loops
			return nil, nil
		default:
			// body
			val, err = n.Body.Eval(ctx)
			if err != nil {
				return nil, err
			}
		}
	}
}
