package javalanche

import (
	"fmt"
	"io"
)

// stateFn is a state function that will return
// the next or nil when tokenizing is finished
type stateFn func(t *Tokenizer) stateFn

// lexText skips whitespace until it's time to start
// assembling a Token
func lexText(t *Tokenizer) stateFn {
	for {
		r, l, err := t.reader.PeekRune()
		switch {
		case err != nil:
			t.emitError(err)
			return nil
		case isDigit(r):
			// number
			return lexNumber
		case isKeywordRune(r):
			// keyword or identifier
			return lexKeyword
		case isIdentifierStart(r):
			// identifier
			return lexIdentifier
		case isDoubleQuote(r):
			// string
			return lexDoubleQuoteString
		case isSingleQuote(r):
			// string
			return lexSingleQuoteString
		case isOperatorStart(r):
			// operator
			return lexOperator
		case isPunctuation(r):
			// punctuation
			return lexPunctuation
		case isWhitespace(r):
			// discard
			t.reader.DiscardBytes(l)
		default:
			err := fmt.Errorf("invalid rune: %q", r)
			t.emitError(err)
			return nil
		}
	}
}

// String ""
func lexDoubleQuoteString(t *Tokenizer) stateFn {
	// open string
	t.acceptFn(isDoubleQuote)
	// inside string
	t.acceptAllFn(func(r rune) bool {
		// match anything that isn't a double quote
		return !isDoubleQuote(r)
	})
	// close string
	closed := t.acceptFn(isDoubleQuote)

	// content
	s := t.reader.EmitString()
	s = s[1:] // remove opening "
	if closed {
		// remove closing "
		s = s[:len(s)-1]
	}

	// emit
	t.emitValue(String, s)
	return lexText
}

// String '
func lexSingleQuoteString(t *Tokenizer) stateFn {
	// open string
	t.acceptFn(isSingleQuote)
	// inside string
	t.acceptAllFn(func(r rune) bool {
		// match anything that isn't a single quote
		return !isSingleQuote(r)
	})
	// close string
	closed := t.acceptFn(isSingleQuote)

	// content
	s := t.reader.EmitString()
	s = s[1:] // remove opening '
	if closed {
		// remove closing '
		s = s[:len(s)-1]
	}

	// emit
	t.emitValue(String, s)
	return lexText
}

// Identify lexing needs
func lexKeyword(t *Tokenizer) stateFn {
	for {
		r, _, err := t.reader.ReadRune()
		switch {
		case err == io.EOF:
			lexEmitKeyword(t)
			return nil
		case err != nil:
			lexEmitKeyword(t)
			t.emitError(err)
			return nil
		case isKeywordRune(r):
			// continue
		case isIdentifierPart(r):
			// continue on lexIndentifier
			return lexIdentifier
		default:
			// end of keyword
			err = t.reader.UnreadRune()
			if err != nil {
				t.emitError(err)
				return nil
			}

			lexEmitKeyword(t)
			return lexText
		}
	}
}

// Lexes Keywords as correct Types
func lexEmitKeyword(t *Tokenizer) {
	s := t.reader.EmitString()
	switch {
	case isBooleanString(s):
		// true or false
		t.emitValue(Boolean, s)
	case isLogicalOperatorString(s):
		// 'and', 'or', 'xor'
		t.emitValue(Operator, s)
	case isKeyword(s):
		// other keywords
		t.emitValue(Keyword, s)
	default:
		// not a keyword
		t.emitValue(Identifier, s)
	}
}

// Lexes Identifier
func lexIdentifier(t *Tokenizer) stateFn {
	for {
		r, _, err := t.reader.ReadRune()
		switch {
		case err == io.EOF:
			t.emitToken(Identifier)
			return nil
		case err != nil:
			t.emitToken(Identifier)
			t.emitError(err)
			return nil
		case isIdentifierPart(r):
			// continue
		default:
			// end of Identifier
			err = t.reader.UnreadRune()
			if err != nil {
				t.emitError(err)
				return nil
			}
			// emit
			t.emitToken(Identifier)
			return lexText
		}
	}
}

// Lexes Number as Int or Float
func lexNumber(t *Tokenizer) stateFn {
	t.acceptAllFn(isDigit)

	if t.accept(".") {
		// dot determines we are emitting float
		t.acceptAllFn(isDigit)
		t.emitToken(Float)
	} else {
		t.emitToken(Integer)
	}

	return lexText
}

// Lexes left and right parenthesis
func lexPunctuation(t *Tokenizer) stateFn {
	// it can't fail because of the previous PeekRune()
	r, _, _ := t.reader.ReadRune()
	switch r {
	case '(':
		t.emitToken(LeftParen)
	case ')':
		t.emitToken(RightParen)
	case '\n':
		t.emitToken(EOL)
	default:
		// can't happen. we know it satisfies isPunctionation()
		panic("unreachable")
	}

	return lexText
}

// Lexes operators and double operators
func lexOperator(t *Tokenizer) stateFn {
	r1, _, _ := t.reader.ReadRune()
	switch {
	case isOperatorWithSecond(r1):
		// these could have a second rune
		r2, _, err := t.reader.ReadRune()
		switch {
		case err != nil:
			// read error, fatal
			t.emitToken(Operator)
			t.emitError(err)
			return nil
		case isOperatorPair(r1, r2):
			// good pair
		case isOperatorNeedsSecond(r1):
			// doesn't work without a second, non-fatal
			err := fmt.Errorf("%q: %s", r1+r2, "invalid operator")
			t.emitError(err)
			return lexText
		default:
			// r2 isn't part of the op, emit without
			if err = t.reader.UnreadRune(); err != nil {
				// unread error, fatal
				t.emitToken(Operator)
				t.emitError(err)
				return nil
			}
		}
	default:
		// single rune op
	}

	// emit and continue
	t.emitToken(Operator)
	return lexText
}
