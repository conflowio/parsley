package text

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"

	"github.com/opsidian/parsley/reader"
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
	return fmt.Sprintf("%d|%d:%d", p.pos, p.line, p.col)
}

// Reader defines a byte reader
type Reader struct {
	b                 []byte
	cur               Position
	charCount         int
	ignoreWhitespaces bool
	regexpCache       map[string]*regexp.Regexp
}

// NewReader creates a new reader instance
func NewReader(b []byte, ignoreWhitespaces bool) *Reader {
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
func (r *Reader) Clone() reader.Reader {
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
			return 0, 0, fmt.Errorf("invalid UTF-8 byte sequence encountered at %s", r.cur)
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

// PeakRune reads the next character but does not move the cursor
func (r *Reader) PeakRune() (ch rune, size int, err error) {
	if r.cur.pos >= len(r.b) {
		return 0, 0, io.EOF
	}
	if c := r.b[r.cur.pos]; c < utf8.RuneSelf {
		ch = rune(c)
		size = 1
	} else {
		ch, size = utf8.DecodeRune(r.b[r.cur.pos:])
		if ch == utf8.RuneError {
			return 0, 0, fmt.Errorf("invalid UTF-8 byte sequence encountered at %s", r.cur)
		}
	}
	return
}

// ReadMatch reads a set of characters matching the given regular expression
func (r *Reader) ReadMatch(expr string, includeWhitespaces bool) ([]string, reader.Position, bool) {
	if r.ignoreWhitespaces && !includeWhitespaces {
		r.readWhitespaces()
	}

	loc := r.getPattern(expr).FindSubmatchIndex(r.b[r.cur.pos:])
	if loc == nil {
		return nil, nil, false
	}
	pos := r.cur
	matches := make([]string, len(loc)/2)
	for i := 0; i < len(loc)/2; i++ {
		matches[i] = string(r.b[r.cur.pos+loc[i*2] : r.cur.pos+loc[i*2+1]])
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

	return matches, pos, true
}

// PeakMatch reads a set of characters matching the given regular expression but doesn't move the cursor
// Also it never ignores whitespaces
func (r *Reader) PeakMatch(expr string) ([]string, bool) {
	pos := r.cur.pos

	loc := r.getPattern(expr).FindSubmatchIndex(r.b[pos:])
	if loc == nil {
		return nil, false
	}

	matches := make([]string, len(loc)/2)
	for i := 0; i < len(loc)/2; i++ {
		matches[i] = string(r.b[pos+loc[i*2] : pos+loc[i*2+1]])
	}

	return matches, true
}

// Readf uses the given function to match the next token
func (r *Reader) Readf(f func(b []byte) (string, int, bool), includeWhitespaces bool) (string, reader.Position, bool) {
	if r.ignoreWhitespaces && !includeWhitespaces {
		r.readWhitespaces()
	}

	pos := r.cur
	value, l, ok := f(r.b[r.cur.pos:])
	if !ok {
		return "", nil, false
	}
	if l != 0 {
		str := string(r.b[r.cur.pos : r.cur.pos+l])
		for _, ch := range str {
			r.charCount--
			if ch != '\n' {
				r.cur.col++
			} else {
				r.cur.line++
				r.cur.col = 1
			}
		}
		r.cur.pos += l
	}
	return value, pos, true
}

// Remaining returns with the remaining character count
func (r *Reader) Remaining() int {
	return r.charCount
}

// Cursor returns with the cursor's position
func (r *Reader) Cursor() reader.Position {
	return r.cur
}

// IsEOF returns true if we reached the end of the buffer
func (r *Reader) IsEOF() bool {
	return r.cur.pos >= len(r.b)
}

func (r *Reader) String() string {
	return fmt.Sprintf("R{%s}", r.cur)
}

func (r *Reader) readWhitespaces() {
	loc := r.getPattern("[ \n\r\t]+").FindIndex(r.b[r.cur.pos:])
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
		rc = regexp.MustCompile("^(?:" + expr + ")")
		r.regexpCache[expr] = rc
	}
	return
}
