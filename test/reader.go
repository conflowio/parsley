// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
