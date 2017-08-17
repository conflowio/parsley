// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package reader

import (
	"fmt"
)

// err is an error with a text position
type err struct {
	cause error
	msg   string
	pos   Position
}

// NewError creates a new reader error instance
func NewError(pos Position, format string, values ...interface{}) Error {
	return err{
		pos: pos,
		msg: fmt.Sprintf(format, values...),
	}
}

// WrapError wraps the given error in a parser error
// If the cause is already a parser error then WrapError returns the same error with an updated error message
func WrapError(pos Position, cause error, format string, values ...interface{}) Error {
	msg := fmt.Sprintf(format, values...)
	if readerErr, ok := cause.(Error); ok {
		pos = readerErr.Pos()
		cause = readerErr.Cause()
		if format == "" {
			msg = readerErr.Msg()
		}
	}
	return err{
		pos:   pos,
		cause: cause,
		msg:   msg,
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
func (e err) Pos() Position {
	return e.pos
}
