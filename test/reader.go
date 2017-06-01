package test

import (
	"fmt"

	"github.com/opsidian/parsley/reader"
)

type Position struct {
	pos int
}

func NewPosition(pos int) Position {
	return Position{pos}
}

func (p Position) Pos() int       { return p.pos }
func (p Position) String() string { return fmt.Sprintf("Pos{%d}", p.pos) }

type Reader struct {
	cur       Position
	remaining int
	isEOF     bool
	cloned    bool
}

func (r Reader) Clone() reader.Reader    { r.cloned = true; return r }
func (r Reader) Cursor() reader.Position { return r.cur }
func (r Reader) Remaining() int          { return r.remaining }
func (r Reader) IsEOF() bool             { return r.isEOF }
func (r Reader) String() string {
	return fmt.Sprintf("Reader{%d %d %v}", r.cur.Pos(), r.remaining, r.isEOF)
}
func NewReader(pos int, remaining int, isEOF bool, cloned bool) Reader {
	return Reader{Position{pos}, remaining, isEOF, cloned}
}
