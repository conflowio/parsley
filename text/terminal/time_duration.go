// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"errors"
	"fmt"
	"time"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// TimeDurationType contains the time duration type's name
const TimeDurationType = "time.Duration"

// TimeDurationNode is a leaf node in the AST
type TimeDurationNode struct {
	value     time.Duration
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewTimeDurationNode creates a new TimeDurationNode instance
func NewTimeDurationNode(value time.Duration, pos parsley.Pos, readerPos parsley.Pos) *TimeDurationNode {
	return &TimeDurationNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (t *TimeDurationNode) Token() string {
	return "TIME_DURATION"
}

// Type returns
func (t *TimeDurationNode) Type() string {
	return TimeDurationType
}

// Value returns with the value of the node
func (t *TimeDurationNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return t.value, nil
}

// Pos returns the position
func (t *TimeDurationNode) Pos() parsley.Pos {
	return t.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (t *TimeDurationNode) ReaderPos() parsley.Pos {
	return t.readerPos
}

// SetReaderPos changes the reader position
func (t *TimeDurationNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	t.readerPos = fun(t.readerPos)
}

// String returns with a string representation of the node
func (t *TimeDurationNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", t.Token(), t.value, t.pos, t.readerPos)
}

// TimeDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func TimeDuration() parser.Func {
	notFoundErr := errors.New("was expecting a valid time duration")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[0-9]+(?:\\.[0-9]+)?(?:ns|us|µs|μs|ms|s|m|h))+"); result != nil {
			duration, err := time.ParseDuration(string(result))
			if err != nil {
				return nil, data.EmptyIntSet, parsley.NewError(pos, err)
			}
			return NewTimeDurationNode(duration, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
