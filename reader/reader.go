package reader

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"
)

// Reader defines a byte reader
type Reader struct {
	b                 []byte
	i                 int
	charCount         int
	line              int
	col               int
	ignoreWhitespaces bool
	regexpCache       map[string]*regexp.Regexp
}

// New creates a new reader instance
func New(b []byte, ignoreWhitespaces bool) *Reader {
	if ignoreWhitespaces {
		b = bytes.Trim(b, "\r\n\t ")
	}
	return &Reader{
		b:                 b,
		i:                 0,
		charCount:         utf8.RuneCount(b),
		line:              1,
		col:               1,
		ignoreWhitespaces: ignoreWhitespaces,
		regexpCache:       make(map[string]*regexp.Regexp),
	}
}

// Clone creates a new reader with the same position
func (r *Reader) Clone() *Reader {
	return &Reader{
		b:                 r.b,
		i:                 r.i,
		charCount:         r.charCount,
		line:              r.line,
		col:               r.col,
		ignoreWhitespaces: r.ignoreWhitespaces,
		regexpCache:       r.regexpCache,
	}
}

// ReadRune reads the next character
func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if r.i >= len(r.b) {
		return 0, 0, io.EOF
	}
	if c := r.b[r.i]; c < utf8.RuneSelf {
		ch = rune(c)
		size = 1
	} else {
		ch, size = utf8.DecodeRune(r.b[r.i:])
		if ch == utf8.RuneError {
			return 0, 0, fmt.Errorf("Invalid UTF-8 byte sequence encountered at %d:%d", r.line, r.col)
		}
	}
	r.i += size
	r.charCount--
	if ch != '\n' {
		r.col++
	} else {
		r.line++
		r.col = 1
	}
	return
}

// ReadMatch reads a set of characters matching the given regular expression
func (r *Reader) ReadMatch(expr string) (matches []string, pos int) {
	if expr[0] != '^' {
		panic("Regexp match should start with ^")
	}

	if r.ignoreWhitespaces {
		r.readWhitespaces()
	}

	loc := r.getPattern(expr).FindSubmatchIndex(r.b[r.i:])
	if loc == nil {
		return nil, -1
	}
	pos = r.i
	matches = make([]string, len(loc)/2)
	matches[0] = string(r.b[r.i : r.i+loc[1]])
	if len(loc) > 2 {
		for i := 1; i < len(loc)/2; i++ {
			matches[i] = string(r.b[r.i+loc[i*2] : r.i+loc[i*2+1]])
		}
	}

	r.i += loc[1]
	for _, ch := range matches[0] {
		r.charCount--
		if ch != '\n' {
			r.col += len(string(ch))
		} else {
			r.line++
			r.col = 1
		}
	}

	return
}

// CharsRemaining returns with the remaining character count
func (r *Reader) CharsRemaining() int {
	return r.charCount
}

// Position returns the current byte index
func (r *Reader) Position() int {
	return r.i
}

// Line returns the current line position
func (r *Reader) Line() int {
	return r.line
}

// Column returns the current column position
func (r *Reader) Column() int {
	return r.col
}

// IsEOF returns true if we reached the end of the buffer
func (r *Reader) IsEOF() bool {
	return r.i >= len(r.b)
}

func (r *Reader) String() string {
	return fmt.Sprintf("Reader{pos: %d, %d chars left}\n", r.i, r.CharsRemaining())
}

func (r *Reader) readWhitespaces() {
	loc := r.getPattern("^[ \n\r\t]+").FindIndex(r.b[r.i:])
	if loc == nil {
		return
	}

	for _, ch := range r.b[r.i : r.i+loc[1]] {
		r.charCount--
		if ch != '\n' {
			r.col++
		} else {
			r.line++
			r.col = 1
		}
	}
	r.i += loc[1]
}

func (r *Reader) getPattern(expr string) (rc *regexp.Regexp) {
	rc, ok := r.regexpCache[expr]
	if !ok {
		rc = regexp.MustCompile(expr)
		r.regexpCache[expr] = rc
	}
	return
}
