package javalanche

import (
	"fmt"
)

// StageNode is either a Token or a Node
type StageNode struct {
	token *Token
	node  Node
}

// Node returns node if one exists
func (n StageNode) Node() (Node, bool) {
	if n.node != nil {
		return n.node, true
	}
	return nil, false
}

// Token returns token if one exists
func (n StageNode) Token() (*Token, bool) {
	if n.token != nil {
		return n.token, true
	}
	return nil, false
}

// Any returns either node or token
func (n StageNode) Any() any {
	switch {
	case n.node != nil:
		return n.node
	default:
		return n.token
	}
}

// Parse returns node or parses a leaf
func (n StageNode) Parse() (Node, error) {
	switch {
	case n.node != nil:
		// ready
		return n.node, nil
	default:
		leaf, err := parseLeaf(n.token)
		if err != nil {
			return nil, err
		}
		return leaf, nil
	}
}

// String forms a string representation of the n.any
func (n StageNode) String() string {
	return fmt.Sprintf("%s", n.Any())
}

// NewStageNode returns stageNode if one exist
func NewStageNode(node Node) (StageNode, bool) {
	if node != nil {
		return StageNode{node: node}, true
	}
	return StageNode{}, false
}

// NewStageToken returns stage Token if one exist
func NewStageToken(token *Token) (StageNode, bool) {
	if token != nil {
		return StageNode{token: token}, true
	}
	return StageNode{}, false
}

// Stage represents stage
type Stage struct {
	nodes []StageNode
}

// Reset sets pairs to 0
func (s *Stage) Reset() {
	s.nodes = s.nodes[:0]
}

// Len provides us with length of stage
func (s Stage) Len() int {
	return len(s.nodes)
}

// IsEmpty checks if all tokens have been parsed
func (s Stage) IsEmpty() bool {
	return len(s.nodes) == 0
}

// replaceRange replaces what's between positions from to until
// with the given Node
func (s *Stage) replaceRange(node Node, from, until int) {
	before := s.nodes[0:from]
	after := s.nodes[until+1:]

	s.Println("replaceRange:", "before:", before)
	s.Println("replaceRange:", "inside:", node)
	s.Println("replaceRange:", "after:", after)

	if n, ok := NewStageNode(node); ok {
		nodes := append(before, n)
		nodes = append(nodes, after...)
		s.nodes = nodes
		return
	}

	panic("unreachable")
}

// AppendTokens appends tokens to stage
func (s *Stage) AppendTokens(tokens ...*Token) {
	for _, token := range tokens {
		switch {
		case isLeafToken(token):
			// convert token of leaf Node immediatelly
			leaf, err := parseLeaf(token)
			switch {
			case err != nil:
				panic("unreachable")
			default:
				s.AppendNodes(leaf)
			}
		default:
			if n, ok := NewStageToken(token); ok {
				s.nodes = append(s.nodes, n)
			}
		}
	}
}

// AppendNodes appends nodes to stage
func (s *Stage) AppendNodes(nodes ...Node) {
	for _, node := range nodes {
		if n, ok := NewStageNode(node); ok {
			s.nodes = append(s.nodes, n)
		}
	}
}

// getNodeBefore returns the Node preceeding the given index
// if it exists
func (s *Stage) getNodeBefore(i int) (Node, error) {
	if i > 0 {
		return s.nodes[i-1].Parse()
	}

	// nothing before the first
	return nil, &ErrInvalidToken{}
}

// getNodeAfter returns the Node immediatelly after the given
// index if it exists
func (s *Stage) getNodeAfter(i int) (Node, error) {
	last := len(s.nodes) - 1
	if i+1 > last {
		// need more data
		return nil, ErrMoreData
	}

	return s.nodes[i+1].Parse()
}

// Parse parses delegates parsing dpeending on len of slice
func (s *Stage) Parse() (Node, error) {
	i := 0
	for {
		s.PrintDetails("Parse: pass:%v", i)
		i++

		switch len(s.nodes) {
		case 0:
			return nil, ErrMoreData
		case 1:
			// single node
			node, err := s.nodes[0].Parse()

			switch {
			case err == ErrMoreData:
				// wait for more data
				return nil, err
			case err != nil:
				// parseError, reset and report
				s.Reset()
				return nil, err
			default:
				// successfully parsed our only pair
				s.Reset()
				return node, nil
			}
		default:
			// brackets
			start, end, found, err := s.findBrackets()
			switch {
			case err != nil:
				// unbalanced or incomplete
				return nil, err
			case found:
				// brackets
				err = s.parseBracketed(start, end)
				if err != nil {
					// bad bracketed
					return nil, err
				}
			default:
				err = s.parseUnbracketed()
				if err != nil {
					// bad unbracketed or incomplete
					return nil, err
				}
			}
		}
	}
}

func (s *Stage) parseBracketed(start, end int) error {
	var result Node

	left := s.nodes[:start]
	right := s.nodes[end+1:]
	inside := s.nodes[start+1 : end]

	s.Println("parseBracketed:", "left:", left)
	s.Println("parseBracketed:", "inside:", inside)
	s.Println("parseBracketed:", "right:", right)

	switch len(inside) {
	case 0:
		// empty bracketed isn't valid
		t, _ := s.nodes[end].Token()
		return &ErrInvalidToken{
			Token:  t,
			Reason: "unexpected ')'",
		}
	case 1:
		// single node
		leaf, err := inside[0].Parse()
		if err != nil {
			return err
		}
		result = leaf
		s.replaceRange(result, start, end)
		return nil
	default:
		// many nodes
		err := s.parseRange(start+1, end)
		if err != nil {
			s.Println("parseBracketed:", "err:", err)
		}
		return nil
	}
}

// parseUnbracketed parses unbracketed nodes, binary and unary
func (s *Stage) parseUnbracketed() error {
	return s.parseRange(0, len(s.nodes))
}

func (s *Stage) parseRange(start, end int) error {
	pivot, op := findHighestPrecedenceOperatorInRange(s.nodes, start, end)
	if op == nil {
		return &ErrInvalidToken{}
	}

	s.Println("parseUnbracketed:", "op:", op, "at", pivot)
	switch {
	case isPrefixUnaryOperator(op.Value):
		// ... op after ...
		after, err := s.getNodeAfter(pivot)
		if err != nil {
			return err
		}

		// could we be on a binary instead?
		if isBinaryOperator(op.Value) {
			before, err := s.getNodeBefore(pivot)
			if err == nil {
				// ... before op after ...
				n := &BinaryExpression{
					Left:  before,
					Op:    op.Value,
					Right: after,
				}

				s.Printf("parseUnbracketed: pivot:%v [%s %s %s] → %s", pivot, before, op, after, n)
				s.replaceRange(n, pivot-1, pivot+1)
				return nil
			}

			// nope, continue as prefixed unary
		}

		n := &UnaryExpression{
			Op:   op.Value,
			Expr: after,
		}

		s.Printf("parseUnbracketed: [%s %s] → %s", op, after, n)
		s.replaceRange(n, pivot, pivot+1)
		return nil
	case isSuffixUnaryOperator(op.Value):
		// ... before op ...
		before, err := s.getNodeBefore(pivot)

		if err != nil {
			return err
		}

		n := &UnaryExpression{
			Expr: before,
			Op:   op.Value,
		}

		s.Printf("parseUnbracketed: [%s %s] → %s", before, op, n)
		s.replaceRange(n, pivot-1, pivot)
		return nil
	case isBinaryOperator(op.Value):
		// ... before op after ...
		before, err := s.getNodeBefore(pivot)
		if err != nil {
			return err
		}

		after, err := s.getNodeAfter(pivot)
		if err != nil {
			return err
		}

		n := &BinaryExpression{
			Left:  before,
			Op:    op.Value,
			Right: after,
		}

		s.Printf("parseUnbracketed: pivot:%v [%s %s %s] → %s", pivot, before, op, after, n)
		s.replaceRange(n, pivot-1, pivot+1)
		return nil

	default:
		return fmt.Errorf("unsupported operator: %s", op.Value)
	}
}

func (s *Stage) findBrackets() (int, int, bool, error) {
	lastOpen := -1

	for i, node := range s.nodes {
		if token, ok := node.Token(); ok {
			switch token.Type {
			case LeftParen:
				lastOpen = i
			case RightParen:
				switch {
				case lastOpen < 0:
					// never opened
					err := &ErrInvalidToken{
						Token:  token,
						Reason: "unmatched closing parenthesis",
					}
					return 0, i, false, err
				default:
					// matched
					return lastOpen, i, true, nil
				}
			}
		}
	}

	switch {
	case lastOpen < 0:
		// none
		return 0, 0, false, nil
	default:
		// never closed, incomplete
		return lastOpen, 0, false, ErrMoreData
	}
}

// ParseLeaf parses leaves
func (s *Stage) ParseLeaf(token *Token) (Node, error) {
	leaf, err := parseLeaf(token)
	switch {
	case err != nil:
		s.Println("ParseLeaf:", token, "→ err:", err)
		return nil, err
	default:
		s.Println("ParseLeaf:", token, "→", leaf)
		return leaf, nil
	}
}

func parseLeaf(token *Token) (Node, error) {
	var leaf Node

	switch token.Type {
	case Identifier:
		leaf = NewVariable(token.Value)
	case Integer:
		leaf, _ = NewIntegerString(token.Value)
	case Float:
		leaf, _ = NewFloatString(token.Value)
	case String:
		leaf = NewString(token.Value)
	case Boolean:
		leaf, _ = NewBooleanString(token.Value)
	}

	switch {
	case leaf == nil:
		return nil, &ErrInvalidToken{token, ""}
	default:
		return leaf, nil
	}
}

func isLeafToken(token *Token) bool {
	if token != nil {
		switch token.Type {
		case Identifier, Integer, Float, String, Boolean:
			return true
		}
	}
	return false
}
