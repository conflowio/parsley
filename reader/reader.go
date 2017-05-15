package reader

import (
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"
)

// Reader defines a byte reader
type Reader struct {
	b                 []byte
	i                 int
	line, col         int
	regexpCache       map[string]*regexp.Regexp
	ignoreWhitespaces bool
}

// New creates a new reader instance
func New(b []byte) *Reader {
	return &Reader{
		b:                 b,
		i:                 0,
		line:              1,
		col:               1,
		regexpCache:       make(map[string]*regexp.Regexp),
		ignoreWhitespaces: true,
	}
}

// Clone creates a new reader with the same position
func (r *Reader) Clone() *Reader {
	return &Reader{
		b:                 r.b,
		i:                 r.i,
		line:              r.line,
		col:               r.col,
		regexpCache:       r.regexpCache,
		ignoreWhitespaces: true,
	}
}

// SetIgnoreWhitespaces sets whether reads should ignore any whitespaces on the left
func (r *Reader) SetIgnoreWhitespaces(ignoreWhitespaces bool) {
	r.ignoreWhitespaces = ignoreWhitespaces
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
	if ch != '\n' {
		r.col++
	} else {
		r.line++
		r.col = 1
	}
	return
}

func (r *Reader) getWhitespaceIndex() int {
	// TODO: column and line is not handled
	loc := r.getPattern("^\\s*").FindIndex(r.b[r.i:])
	if loc == nil {
		return 0
	}
	return loc[1]
}

// ReadMatch reads a set of characters matching the given regular expression
func (r *Reader) ReadMatch(expr string) (matches []string, pos int) {
	if expr[0] != '^' {
		panic("Regexp match should start with ^")
	}

	whitespaceIndex := 0
	if r.ignoreWhitespaces {
		whitespaceIndex = r.getWhitespaceIndex()
	}

	loc := r.getPattern(expr).FindSubmatchIndex(r.b[r.i+whitespaceIndex:])
	if loc == nil {
		return nil, -1
	}
	r.i += whitespaceIndex
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
		if ch != '\n' {
			r.col += len(string(ch))
		} else {
			r.line++
			r.col = 1
		}
	}

	return
}

// Position returns the current byte index
func (r *Reader) Position() int {
	return r.i
}

// TestPosition returns the current line and column position
func (r *Reader) TestPosition() (int, int) {
	return r.line, r.col
}

func (r *Reader) getPattern(expr string) (rc *regexp.Regexp) {
	rc, ok := r.regexpCache[expr]
	if !ok {
		rc = regexp.MustCompile(expr)
		r.regexpCache[expr] = rc
	}
	return
}

// ReadEOF returns true if we reached the end of the buffer
func (r *Reader) ReadEOF() bool {
	if r.ignoreWhitespaces {
		r.i += r.getWhitespaceIndex()
	}
	return r.i >= len(r.b)
}
