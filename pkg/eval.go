package javalanche

import (
	"errors"
	"strings"
)

func EvalString(expr string) (Value, error) {
	if expr == "" {
		return nil, errors.New("empty string received")
	}
	// Since Tokenizer implements idea of I/O we need to turn our input into the reader
	reader := strings.NewReader(expr)
	lexer := NewTokenizer(reader)
	evaluator := NewEvaluator()
	parser := NewParser(lexer, evaluator)

	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	value, err := ast.Eval()
	if err != nil {
		return nil, err
	}

	return value, nil
}
