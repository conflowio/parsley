// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"errors"
	"fmt"
	"strings"
)

// Error is an error with a position
type Error struct {
	err error
	pos Position
}

// NewError creates a new error with the given position
// If the passed error is already a parsley.Error it returns the original error
// as it should have already the correct position.
func NewError(err error, pos Position) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{
		err: err,
		pos: pos,
	}
}

// Cause returns with the original error
func (e *Error) Cause() error {
	return e.err
}

// Error returns with the full error message including the position
func (e *Error) Error() string {
	if e.pos == nil {
		return e.err.Error()
	}

	if e.pos == NilPosition {
		return e.err.Error()
	}

	return fmt.Sprintf("%s at %s", e.err.Error(), e.pos.String())
}

// Pos returns with the error's position
func (e *Error) Pos() Position {
	return e.pos
}

// WrapError wraps the given error in a error
// If format contains the "{{err}}" placeholder it will be replaced with the original error message
func WrapError(err *Error, format string, values ...interface{}) *Error {
	msg := fmt.Sprintf(format, values...)
	if msg == "" {
		return err
	}
	return &Error{
		err: errors.New(strings.Replace(msg, "{{err}}", err.Cause().Error(), -1)),
		pos: err.Pos(),
	}
}
