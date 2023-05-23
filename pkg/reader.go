package javalanche

import (
	"errors"
	"io"
	"unicode/utf8"
)

const (
	// ReadBufferSize indicates the
	// initial buffer size
	ReadBufferSize = 2
)

// Implemented Interfaces
var (
	_ io.RuneScanner = (*Reader)(nil)
)

// Reader is a custom RuneScanner for composing tokens
type Reader struct {
	src io.Reader
	buf []byte

	cursor       int
	lastRune     rune
	lastRuneSize int
}

// NewReader builds a lexer Reader on top of a regular io.Reader
func NewReader(rd io.Reader) *Reader {
	if rd == nil {
		return nil
	}

	return &Reader{
		buf: make([]byte, 0, ReadBufferSize),
		src: rd,
	}
}

func (b *Reader) fill(needed int) error {
	for len(b.buf) < needed {
		// slice with length needed
		newBuf := make([]byte, needed)
		// Copy existing data into the new buffer
		copy(newBuf, b.buf)

		// Allocate additional space in the slice for the read
		readBuf := newBuf[len(b.buf):needed]

		// Read from the source into the newly allocated space
		n, err := b.src.Read(readBuf)

		// return error other than eof
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		// Only use what we need
		newBuf = newBuf[:len(b.buf)+n]
		b.buf = newBuf

		// If we hit EOF, stop fill
		if errors.Is(err, io.EOF) {
			break
		}
	}

	return nil
}

func (b *Reader) needsBytes(n int) error {
	if n < 1 {
		n = 1
	}

	return b.fill(b.cursor + n)
}

// ReadRune reads the next rune
func (b *Reader) ReadRune() (rune, int, error) {
	// we need at least one byte to start
	count := 1
	for {
		err := b.needsBytes(count)
		if err != nil {
			return 0, 0, err
		}

		if utf8.FullRune(b.buf[b.cursor:]) {
			// we have a full rune
			break
		}

		// we need more bytes
		count++
	}

	// decode rune
	r, l := utf8.DecodeRune(b.buf[b.cursor:])

	// move the cursor after the decoded rune
	b.cursor += l

	// remember result for UnreadRune before returning
	b.lastRune = r
	b.lastRuneSize = l

	return r, l, nil
}

// UnreadRune puts back the rune from last ReadRune
func (b *Reader) UnreadRune() error {
	if b.lastRuneSize > 0 {
		// put last rune back into the buffer
		cursor := b.cursor - b.lastRuneSize

		if cursor < 0 {
			// bad new position
			panic("inconsistent buffer")
		}

		// unread
		b.cursor = cursor
		// and make sure we don't unread it again
		b.lastRune = 0
		b.lastRuneSize = 0

		return nil
	}

	return errors.New("invalid UnreadRune call")
}

// PeekRune returns the next rune without moved the cursor
func (b *Reader) PeekRune() (rune, int, error) {
	r, l, err := b.ReadRune()
	if err == nil {
		err = b.UnreadRune()
	}
	return r, l, err
}

// Discard discards everything up to the cursor
func (b *Reader) Discard() error {
	return b.DiscardBytes(0)
}

// DiscardBytes discards everything up to the cursor
// plus some extra bytes
func (b *Reader) DiscardBytes(count int) error {
	switch {
	case b.cursor <= 0 && count == 0:
		// Nothing to do
		return nil
	case count < 0 || count > len(b.buf)-b.cursor:
		// count is invalid
		return errors.New("invalid skip count")
	default:
		// discard everything before b.cursor+count
		copy(b.buf, b.buf[b.cursor+count:])
		b.buf = b.buf[:len(b.buf)-b.cursor-count]
		b.cursor = 0
		// and forget lost runes
		b.lastRune = 0
		b.lastRuneSize = 0
		return nil
	}
}

// Emit returns a byte array corresponding
// to all accepted runes, and resets the buffer
// afterwards
func (b *Reader) Emit() []byte {
	s := make([]byte, b.cursor)
	copy(s, b.buf[:b.cursor])
	b.DiscardBytes(0)
	return s
}

// EmitString returns a string corresponding
// to all accepted runes, and resets the buffer
// afterwards
func (b *Reader) EmitString() string {
	return string(b.Emit())
}

// Accept consumes a rune if it's accepted by the provided
// match function
func (b *Reader) Accept(match func(rune) bool) bool {
	r, _, err := b.ReadRune()
	switch {
	case err != nil:
		// read failed
		return false
	case match != nil && match(r):
		// accepted
		return true
	default:
		// not accepted
		_ = b.UnreadRune()
		return false
	}
}

// AcceptAll consumes a series of runes accepted by the
// provided match function
func (b *Reader) AcceptAll(match func(rune) bool) bool {
	var consumed bool

	for {
		if !b.Accept(match) {
			return consumed
		}
		consumed = true
	}
}
