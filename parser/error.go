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
	Cause() error
	Msg() string
	Pos() reader.Position
}

// err is a parser error
type err struct {
	cause error
	msg   string
	pos   reader.Position
}

// NewError creates a new parser error instance
func NewError(pos reader.Position, format string, values ...interface{}) Error {
	return err{
		pos: pos,
		msg: fmt.Sprintf(format, values...),
	}
}

// WrapError wraps the given error in a parser error
// If the cause is already a parser error then WrapError returns the same error with an updated error message
func WrapError(pos reader.Position, cause error, format string, values ...interface{}) Error {
	if parserErr, ok := cause.(Error); ok {
		pos = parserErr.Pos()
		cause = parserErr.Cause()
	}
	return err{
		pos:   pos,
		cause: cause,
		msg:   fmt.Sprintf(format, values...),
	}
}

// Cause returns with the original error
func (e err) Cause() error {
	return e.cause
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
