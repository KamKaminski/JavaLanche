package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("reader.txt")
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inputs := []string{}

	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	for i, input := range inputs {
		fmt.Printf("Test case #%d: %s\n", i+1, input)
		reader := strings.NewReader(input)
		tokenizer := NewTokenizer(reader)
		parser := NewParser(tokenizer)
		node, err := parser.Parse()

		if err != nil {
			fmt.Printf("Error parsing input: %s\n", err)
			continue
		}

		result, err := node.Eval()
		if err != nil {
			fmt.Printf("Error evaluating expression: %s\n", err)
			continue
		}

		if result.Type() == ValueTypeBool {
			fmt.Printf("Result: %t\n", result.AsBool())
		} else {
			if float64(int(result.AsFloat64())) == result.AsFloat64() {
				fmt.Printf("Result: %d\n", int(result.AsFloat64()))
			} else {
				fmt.Printf("Result: %.2f\n", result.AsFloat64())
			}
		}
	}
}
