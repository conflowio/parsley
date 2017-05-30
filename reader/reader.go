package reader

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"
)

// Position represents a token position
type Position struct {
	pos  int
	line int
	col  int
}

// NewPosition creates a new position instance
func NewPosition(pos int, line int, col int) Position {
	return Position{pos, line, col}
}

// Pos returns with the byte position
func (p Position) Pos() int {
	return p.pos
}

// Line returns with the line position
func (p Position) Line() int {
	return p.line
}

// Col returns with the column position
func (p Position) Col() int {
	return p.col
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.col)
}

// Reader defines a byte reader
type Reader struct {
	b                 []byte
	cur               Position
	charCount         int
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
		cur:               NewPosition(0, 1, 1),
		charCount:         utf8.RuneCount(b),
		ignoreWhitespaces: ignoreWhitespaces,
		regexpCache:       make(map[string]*regexp.Regexp),
	}
}

// Clone creates a new reader with the same position
func (r *Reader) Clone() *Reader {
	return &Reader{
		b:                 r.b,
		cur:               r.cur,
		charCount:         r.charCount,
		ignoreWhitespaces: r.ignoreWhitespaces,
		regexpCache:       r.regexpCache,
	}
}

// ReadRune reads the next character
func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if r.cur.pos >= len(r.b) {
		return 0, 0, io.EOF
	}
	if c := r.b[r.cur.pos]; c < utf8.RuneSelf {
		ch = rune(c)
		size = 1
	} else {
		ch, size = utf8.DecodeRune(r.b[r.cur.pos:])
		if ch == utf8.RuneError {
			return 0, 0, fmt.Errorf("Invalid UTF-8 byte sequence encountered at %s", r.cur)
		}
	}
	r.cur.pos += size
	r.charCount--
	if ch != '\n' {
		r.cur.col++
	} else {
		r.cur.line++
		r.cur.col = 1
	}
	return
}

// ReadMatch reads a set of characters matching the given regular expression
func (r *Reader) ReadMatch(expr string) (matches []string, pos Position) {
	if expr[0] != '^' {
		panic("Regexp match should start with ^")
	}

	if r.ignoreWhitespaces {
		r.readWhitespaces()
	}

	loc := r.getPattern(expr).FindSubmatchIndex(r.b[r.cur.pos:])
	if loc == nil {
		return nil, Position{-1, -1, -1}
	}
	pos = r.cur
	matches = make([]string, len(loc)/2)
	matches[0] = string(r.b[r.cur.pos : r.cur.pos+loc[1]])
	if len(loc) > 2 {
		for i := 1; i < len(loc)/2; i++ {
			matches[i] = string(r.b[r.cur.pos+loc[i*2] : r.cur.pos+loc[i*2+1]])
		}
	}

	r.cur.pos += loc[1]
	for _, ch := range matches[0] {
		r.charCount--
		if ch != '\n' {
			r.cur.col++
		} else {
			r.cur.line++
			r.cur.col = 1
		}
	}

	return
}

// CharsRemaining returns with the remaining character count
func (r *Reader) CharsRemaining() int {
	return r.charCount
}

// Cursor returns with the cursor's position
func (r *Reader) Cursor() Position {
	return r.cur
}

// WithCursor clones the readers and sets the cursor to the specified position
func (r *Reader) WithCursor(pos int, line int, col int) *Reader {
	r2 := r.Clone()
	r2.cur = NewPosition(pos, line, col)
	r2.charCount = utf8.RuneCount(r2.b[r2.cur.Pos():])
	return r2
}

// IsEOF returns true if we reached the end of the buffer
func (r *Reader) IsEOF() bool {
	return r.cur.pos >= len(r.b)
}

func (r *Reader) String() string {
	return fmt.Sprintf("R{%s}", r.cur)
}

func (r *Reader) readWhitespaces() {
	loc := r.getPattern("^[ \n\r\t]+").FindIndex(r.b[r.cur.pos:])
	if loc == nil {
		return
	}

	for _, ch := range r.b[r.cur.pos : r.cur.pos+loc[1]] {
		r.charCount--
		if ch != '\n' {
			r.cur.col++
		} else {
			r.cur.line++
			r.cur.col = 1
		}
	}
	r.cur.pos += loc[1]
}

func (r *Reader) getPattern(expr string) (rc *regexp.Regexp) {
	rc, ok := r.regexpCache[expr]
	if !ok {
		rc = regexp.MustCompile(expr)
		r.regexpCache[expr] = rc
	}
	return
}
