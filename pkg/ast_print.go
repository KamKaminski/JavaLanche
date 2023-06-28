package javalanche

import "fmt"

// PrintNodeHandler is callback ref
type PrintNodeHandler func(...Value) error

// PrintNode represents a node that prints the values of other nodes
type PrintNode struct {
	Nodes   []Node
	Handler PrintNodeHandler
}

// DefaultPrintHandler makes prints in human friendly way
func DefaultPrintHandler(values ...Value) error {
	a := make([]any, 0, len(values))
	for _, v := range values {
		if v != nil {
			a = append(a, v)
		}
	}

	if len(a) > 0 {
		fmt.Println(a...)
	}

	return nil
}

// Eval implements the Node interface
func (n *PrintNode) Eval(ctx *Javalanche) (Value, error) {
	var values []Value

	for _, n := range n.Nodes {
		val, err := n.Eval(ctx)
		switch {
		case err != nil:
			return nil, err
		case val != nil:
			values = append(values, val)
		}
	}

	switch {
	case n.Handler == nil:
		DefaultPrintHandler(values...)
		return nil, nil
	default:
		err := n.Handler(values...)
		return nil, err
	}
}

// AppendNodes appends nodes to the PrintNode's list of nodes
func (n *PrintNode) AppendNodes(nodes ...Node) {
	n.Nodes = append(n.Nodes, nodes...)
}

// NewPrintNode returns a new PrintNode instance
func NewPrintNode(handler PrintNodeHandler, nodes ...Node) *PrintNode {
	return &PrintNode{
		Nodes:   nodes,
		Handler: handler,
	}
}
