package javalanche

import (
	"strings"
	"unicode"
)

// Alphabet used within JaVaLanche lexer
const (
	asciiLetterRunes        = "abcdefghijklmnopqrstuvwxyz"
	operatorWithSecondRunes = "&|=<>!+-"
	operatorStartRunes      = operatorWithSecondRunes + "+-*/%:^"
	punctuationRunes        = "()\n"
)

// isKeywordRune checks if a given rune is a part of ASCII letter runes,
func isKeywordRune(r rune) bool {
	return strings.ContainsRune(asciiLetterRunes, r)
}

// isIdentifierStart is a function that check if are dealing with start of an identifier,
func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// isIdentifierStart is a function that check if are dealing with an identifier,
func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || unicode.IsDigit(r)
}

// isWhiteSpace checks for whitespace
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// isDigit checks for digits
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// isDoubleQuote checks for strings starting with double quotes
func isDoubleQuote(r rune) bool {
	return r == '"'
}

// isSingleQuote checks for strings starting with single quotes
func isSingleQuote(r rune) bool {
	return r == '\''
}

// isOperatorStart handles operator recognition
func isOperatorStart(r rune) bool {
	return strings.ContainsRune(operatorStartRunes, r)
}

// isOperatorWithSecond handles operator recognition
func isOperatorWithSecond(r rune) bool {
	return strings.ContainsRune(operatorWithSecondRunes, r)
}

// isOperatorNeedsSecond handles operator recognition
func isOperatorNeedsSecond(r rune) bool {
	switch r {
	case '&', '|':
		return true
	default:
		return false
	}
}

// isOperatorPair handles double operator recognition
func isOperatorPair(r1 rune, r2 rune) bool {
	op := string([]rune{r1, r2})
	switch op {
	case "==", "<=", ">=", "!=":
		return true
	case "&&", "||":
		return true
	case "++", "--":
		return true
	default:
		return false
	}
}

// isPunctuation recognies parantheses
func isPunctuation(r rune) bool {
	return strings.ContainsRune(punctuationRunes, r)
}

// isBooleanString identifies booleans
func isBooleanString(code string) bool {
	switch code {
	case "true", "false":
		return true
	default:
		return false
	}
}

// isLogicalOperatorStrings handles logical operator detection
func isLogicalOperatorString(code string) bool {
	switch code {
	case "and", "or", "xor":
		return true
	default:
		return false
	}
}

// isKeyword checks if keyword was detected
func isKeyword(code string) bool {
	for _, keyword := range keywords {
		if keyword == code {
			return true
		}
	}
	return false
}

// isBinaryOperator checks if the strings is a binary operator
func isBinaryOperator(code string) bool {
	switch code {
	case "+", "-", "*", "/", "==", "!=", ">", "<", ">=", "<=", "&&", "||", "^", "=", "and", "or":
		return true
	default:
		return false
	}
}

// isPrefixUnaryOperator checks if the strings is a prefix unary operator
func isPrefixUnaryOperator(code string) bool {
	switch code {
	case "!":
		return true
	default:
		return false
	}
}

// isSuffixUnaryOperator checks if the strings is a suffix unary operator
func isSuffixUnaryOperator(code string) bool {
	switch code {
	case "++", "--":
		return true
	default:
		return false
	}
}
