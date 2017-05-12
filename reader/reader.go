package reader

import (
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"
)

// Reader defines a byte reader
type Reader struct {
	b           []byte
	i           int64
	line, col   int
	regexpCache map[string]*regexp.Regexp
}

// New creates a new reader instance
func New(b []byte) *Reader {
	return &Reader{
		b:           b,
		line:        1,
		col:         1,
		regexpCache: make(map[string]*regexp.Regexp),
	}
}

// Clone creates a new reader with the same position
func (r *Reader) Clone() *Reader {
	return &Reader{
		b:           r.b,
		i:           r.i,
		line:        r.line,
		col:         r.col,
		regexpCache: r.regexpCache,
	}
}

// ReadRune reads the next character
func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if r.i >= int64(len(r.b)) {
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
	r.i += int64(size)
	if ch != '\n' {
		r.col++
	} else {
		r.line++
		r.col = 1
	}
	return
}

// ReadMatch reads a set of characters matching the given regular expression
func (r *Reader) ReadMatch(expr string) (str string, size int, ok bool) {
	if expr[0] != '^' {
		panic("Regexp match should start with ^")
	}
	loc := r.getPattern(expr).FindIndex(r.b[r.i:])
	if loc == nil {
		return "", 0, false
	}
	if loc[0] != 0 {
		// This should never happen
		panic("Regexp expression should match the beginning of the text")
	}
	str = string(r.b[r.i : r.i+int64(loc[1])])
	size = loc[1] - loc[0]
	r.i += int64(size)
	for _, ch := range str {
		if ch != '\n' {
			r.col += len(string(ch))
		} else {
			r.line++
			r.col = 1
		}
	}
	return str, size, true
}

// Position returns the current line and column position
func (r *Reader) Position() (int, int) {
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
