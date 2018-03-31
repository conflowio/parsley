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
type Error interface {
	error
	Cause() error
	Pos() Position
}

type err struct {
	cause error
	pos   Position
}

// NewError creates a new error with the given position
// If the passed error is already a parsley.Error it returns the original error
// as it should have already the correct position.
func NewError(cause error, pos Position) Error {
	if e, ok := cause.(*err); ok {
		return e
	}
	return &err{
		cause: cause,
		pos:   pos,
	}
}

// Cause returns with the original error
func (e *err) Cause() error {
	return e.cause
}

// Error returns with the full error message including the position
func (e *err) Error() string {
	if e.pos == nil {
		return e.cause.Error()
	}

	if e.pos == NilPosition {
		return e.cause.Error()
	}

	return fmt.Sprintf("%s at %s", e.cause.Error(), e.pos.String())
}

// Pos returns with the error's position
func (e *err) Pos() Position {
	return e.pos
}

// WrapError wraps the given error in a error
// If format contains the "{{err}}" placeholder it will be replaced with the original error message
func WrapError(e Error, format string, values ...interface{}) Error {
	msg := fmt.Sprintf(format, values...)
	if msg == "" {
		return e
	}
	return &err{
		cause: errors.New(strings.Replace(msg, "{{err}}", e.Cause().Error(), -1)),
		pos:   e.Pos(),
	}
}
