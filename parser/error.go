// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"

	"github.com/opsidian/parsley/reader"
)

// Error is a parser error interface
type Error interface {
	error
	Msg() string
	Pos() reader.Position
}

// err is a parser error
type err struct {
	msg string
	pos reader.Position
}

// NewError creates a new parser error instance
func NewError(pos reader.Position, format string, values ...interface{}) Error {
	return err{
		msg: fmt.Sprintf(format, values...),
		pos: pos,
	}
}

// Msg returns with the error message
func (e err) Msg() string {
	return e.msg
}

// Error returns with the full error message including the position
func (e err) Error() string {
	return fmt.Sprintf("%s at %s", e.msg, e.pos)
}

// Pos returns with the error's position
func (e err) Pos() reader.Position {
	return e.pos
}
