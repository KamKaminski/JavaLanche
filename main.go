// Package main implements REPL
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"

	javalanche "javalanche/pkg"
)

var (
	purple = color.New(color.FgHiYellow).SprintFunc()
	red    = color.New(color.FgHiRed).SprintFunc()
	yellow = color.New(color.FgHiGreen).SprintFunc()
)

// getline reads a line of text from the console, optionally showing a prompt first
func getline(prompt bool) (string, error) {
	if prompt {
		fmt.Print(yellow("javalanche> ")) // Print a prompt
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Trim the newline character from the input
	line = strings.TrimSpace(line)
	return line, nil
}

// printError represents function printing errors
func printError(err error) {
	fmt.Println(red("Error"), err)
}

// printResults prints a message on the console corresponding to the result value of evaluating and expression
func printResult(result javalanche.Value) {
	fmt.Printf("%s %v\n", purple("Result:"), result)
}

// isExit checks if the input is the exit command
func isExit(line string) bool {
	return strings.EqualFold(line, "exit")
}

// Main creates new context
func main() {
	ctx := javalanche.New()

	if isatty.IsTerminal(os.Stdin.Fd()) {
		repl(ctx)
	} else {
		replNonTTY(ctx)
	}
}

// Repl allows user to interact with the program
func repl(ctx *javalanche.Javalanche) {
	prompt := true

	for {
		line, err := getline(prompt)
		switch {
		case err == io.EOF || isExit(line):
			return
		case err != nil:
			printError(err)
			return
		}

		// evaluate
		value, err := ctx.EvalLine(line)
		switch {
		case err == javalanche.ErrMoreData:
			// needs another line
			prompt = false
		case err != nil:
			// error
			printError(err)
			prompt = true
		case value != nil:
			// result
			printResult(value)
			prompt = true
		default:
			// silent statement
			prompt = true
		}
	}
}

// replNonTTY spawns console in a case of main console being non tty
func replNonTTY(ctx *javalanche.Javalanche) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		value, err := ctx.EvalLine(line)
		if err != nil {
			printError(err)
		} else if value != nil {
			printResult(value)
		}
	}
	if err := scanner.Err(); err != nil {
		printError(err)
	}
}
