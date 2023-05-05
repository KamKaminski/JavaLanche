package javalanche

import (
	"errors"
	"fmt"
	"strings"
)

type Evaluator struct {
	Variable map[string]Value
}

// creates a new Evaluator with an empty variable map
func NewEvaluator() *Evaluator {
	return &Evaluator{
		Variable: make(map[string]Value),
	}
}

func (e *Evaluator) SetValue(name string, v Value) error {
	e.Variable[name] = v
	return nil
}

func (e *Evaluator) GetValue(name string) (Value, error) {
	if v, ok := e.Variable[name]; ok {
		return v, nil

	}
	return nil, fmt.Errorf("variable %q not found", name)
}

func EvalString(expr string, evaluator *Evaluator) (Value, error) {
	if expr == "" {
		return nil, errors.New("empty string received")
	}

	// Since Tokenizer implements the idea of I/O, we need to turn our input into the reader
	reader := strings.NewReader(expr)
	lexer := NewTokenizer(reader)
	parser := NewParser(lexer, evaluator)

	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	value, err := ast.Eval(evaluator)
	fmt.Printf("Generated AST: %v\n", ast)
	fmt.Printf("Evaluator state: %v\n", evaluator.Variable)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Evaluator state: %v\n", evaluator.Variable)
	return value, nil
}
