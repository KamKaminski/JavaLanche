package javalanche

import (
	"fmt"
	"strings"
)

// TokenType represents type of our tokens
type TokenType int

// TokenTypes
const (
	Unknown TokenType = iota
	Keyword
	Identifier
	Integer
	Float
	Operator
	Separator
	Boolean
	String
	LeftParen
	RightParen
	EOL
	EOF
)

// TokenTypes as string
func (t TokenType) String() string {
	switch t {
	case Keyword:
		return "Keyword"
	case Identifier:
		return "Identifier"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Operator:
		return "Operator"
	case Separator:
		return "Separator"
	case Boolean:
		return "Boolean"
	case String:
		return "String"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case EOL:
		return "EOL"
	default:
		return "Unknown"
	}
}

// Token represents our Token
type Token struct {
	Type  TokenType
	Value string
}

// GoString  does a recursive print
func (t Token) GoString() string {
	s := t.Type.String()
	switch {
	case s == "Unknown":
		return fmt.Sprintf("&%T{%v, %q}",
			t, int(t.Type), t.Value)
	default:
		return fmt.Sprintf("&%T{%s, %q}",
			t, s, t.Value)
	}
}

// IsType checks if the token has the specified type
// and offers special
func (t Token) IsType(typ TokenType) bool {
	switch {
	case t.Type == typ:
		return true
	case typ == Unknown:
		return t.Type < Keyword || t.Type > EOF
	default:
		return false
	}
}

// Is checks if the token has the specified type and
// optionally one of the provided
func (t *Token) Is(typ TokenType, values ...string) bool {
	switch {
	case t == nil:
		// no Token
		return false
	case !t.IsType(typ):
		// wrong type
		return false
	case len(values) == 0:
		// right type, any value
		return true
	default:
		// any of the acceptable values
		for _, value := range values {
			if t.Value == value {
				return true
			}
		}
		return false
	}
}

// TokenResult Represents result of the token
type TokenResult struct {
	Token *Token
	Err   error
}

// GoString prints Recurisvely  result of tokenisation
func (r TokenResult) GoString() string {
	var buf strings.Builder
	var s string
	if r.Err == nil {
		s = "nil"
	} else {
		s = fmt.Sprintf("%#v", r.Err)
	}

	fmt.Fprintf(&buf, "&%T{%#v, %s}", r, r.Token, s)
	return buf.String()
}
