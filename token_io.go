package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type TokenType int

const (
	Unknown TokenType = iota
	Keyword
	Identifier
	Number
	Operator
	Separator
	Boolean
	String
	Function
	Comment
	LeftParen
	RightParen
	EOF
)

var keywords = []string{
	"if",
	"else",
	"for",
	"while",
	"funx",
	"return",
	"var",
	"true",
	"false",
	"int",
	"float",
	"bool",
	"string",
	"const",
	"print",
	"in",
}

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	reader   *bufio.Reader
	writer   io.Writer
	position int
	buffer   []Token
}

func isLeftParen(r rune) bool {
	return r == '('
}

func isRightParen(r rune) bool {
	return r == ')'
}
func isFunctionLiteral(code string) bool {
	return strings.HasPrefix(code, "funx")
}

func isBooleanLiteral(code string) bool {
	return strings.HasPrefix(code, "true") || strings.HasPrefix(code, "false")
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isDot(r rune) bool {
	return r == '.'
}

func isOperator(code string) bool {
	return strings.Contains("+-*/%=:<>!&|^", code) || code == "and" || code == "or"
}

func isSeparator(r rune) bool {
	return strings.ContainsRune("{}[],;.", r)
}

func isStringDelimiter(r rune) bool {
	return r == '"' || r == '\''
}

func isComment(code string) bool {
	return strings.HasPrefix(code, "#")
}

func isKeyword(code string) bool {
	for _, keyword := range keywords {
		if keyword == code {
			return true
		}
	}
	return false
}

func classifyRune(r rune, code string) TokenType {
	switch {
	case isFunctionLiteral(code):
		return Function
	case isBooleanLiteral(code):
		return Boolean
	case isIdentifierStart(r):
		return Identifier
	case isDigit(r):
		return Number
	case isOperator(code):
		return Operator
	case isLeftParen(r):
		return LeftParen
	case isRightParen(r):
		return RightParen
	case isSeparator(r):
		return Separator
	case isStringDelimiter(r):
		return String
	case isComment(code):
		return Comment
	default:
		return Unknown
	}
}

func isValidMultiCharOperator(code string) bool {
	switch code {
	case "==", "!=", "<=", ">=", "&&", "||":
		return true
	default:
		return false
	}
}
func (t *Tokenizer) PeekToken() (*Token, error) {
	token, err := t.NextToken()
	if err != nil {
		return nil, err
	}
	t.PushBack(token)
	return token, nil
}

func readTestCasesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var testCases []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testCases = append(testCases, scanner.Text())
	}

	return testCases, scanner.Err()
}
func (t TokenType) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Keyword:
		return "Keyword"
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	case Operator:
		return "Operator"
	case Separator:
		return "Separator"
	case Boolean:
		return "Boolean"
	case String:
		return "String"
	case Function:
		return "Function"
	case Comment:
		return "Comment"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	default:
		return "Unknown"
	}
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader:   bufio.NewReader(r),
		position: 0,
	}
}
func (t *Tokenizer) PushBack(token *Token) {
	t.buffer = append(t.buffer, *token)
}

func (t *Tokenizer) Run() {
	for {
		token, err := t.NextToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(t.writer, "Error: %v\n", err)
			break
		}

		fmt.Fprintf(t.writer, "Type: %v, Value: %v\n", token.Type.String(), token.Value)
	}
}

func (t *Tokenizer) NextToken() (*Token, error) {
	var current strings.Builder

	for {
		r, _, err := t.reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		t.position++
		code := string(r)
		tokenType := classifyRune(r, code)
		switch tokenType {
		case Identifier:
			for isIdentifierPart(r) {
				current.WriteRune(r)
				r, _, err = t.reader.ReadRune()
				if err != nil {
					break
				}
			}
			t.reader.UnreadRune()

			value := current.String()
			if isKeyword(value) {
				tokenType = Keyword
			} else if value == "and" || value == "or" {
				tokenType = Operator
			}
			return &Token{Type: tokenType, Value: value}, nil

		case Number:
			for isDigit(r) || isDot(r) {
				current.WriteRune(r)
				r, _, err = t.reader.ReadRune()
				if err != nil {
					break
				}
			}
			t.reader.UnreadRune()

			return &Token{Type: tokenType, Value: current.String()}, nil

		case Boolean, Function:
			return &Token{Type: tokenType, Value: code}, nil

		case Operator:
			current.WriteRune(r)
			nextR, _, err := t.reader.ReadRune()
			if err == nil {
				temp := current.String() + string(nextR)
				if isValidMultiCharOperator(temp) {
					current.WriteRune(nextR)
				} else {
					t.reader.UnreadRune()
				}
			} else {
				t.reader.UnreadRune()
			}
			return &Token{Type: tokenType, Value: current.String()}, nil

		case Separator:
			return &Token{Type: tokenType, Value: string(r)}, nil
		case String:
			delimiter := r
			current.WriteRune(delimiter)
			code = ""
			for {
				r, _, err = t.reader.ReadRune()
				if err != nil {
					break
				}
				code += string(r)

				if r == delimiter {
					break
				}
			}
			current.WriteString(code)
			return &Token{Type: tokenType, Value: current.String()}, nil

		case Comment:
			current.WriteRune(r)
			for {
				r, _, err = t.reader.ReadRune()
				if err != nil || r == '\n' {
					break
				}
				current.WriteRune(r)
			}
			if err != nil {
				t.reader.UnreadRune()
			}
			return &Token{Type: tokenType, Value: current.String()}, nil

		case LeftParen, RightParen:
			current.WriteRune(r)
			return &Token{Type: tokenType, Value: current.String()}, nil

		case Unknown:
			if isWhitespace(r) {
				code = code[1:]
			} else {
				// Report the correct position using the `position` variable
				return nil, fmt.Errorf("Invalid character at position %d in remaining code: %s", t.position, code)
			}
		}
	}
	return nil, io.EOF

}
