// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"fmt"
	"strings"
)

// Error is an error with a position
//go:generate counterfeiter . Error
type Error interface {
	Error() string
	Cause() error
	Pos() Pos
}

type err struct {
	cause error
	msg   string
	pos   Pos
}

// NewError creates a new error with the given position
// If the passed error is already a parsley.Error it returns the original error
// as it should have already the correct position.
func NewError(pos Pos, cause error) Error {
	switch c := cause.(type) {
	case Error:
		return c
	default:
		return &err{
			cause: cause,
			msg:   cause.Error(),
			pos:   pos,
		}
	}
}

// NewErrorf creates a new error with the given position and message
func NewErrorf(pos Pos, format string, values ...interface{}) Error {
	cause := fmt.Errorf(format, values...)
	return &err{
		cause: cause,
		msg:   cause.Error(),
		pos:   pos,
	}
}

// Error returns with the full error message including the position
func (e *err) Error() string {
	return e.msg
}

// Pos returns with the error's position
func (e *err) Pos() Pos {
	return e.pos
}

// Cause returns with the original error
func (e *err) Cause() error {
	return e.cause
}

// WrapError wraps the given error in a error
// If format contains the "{{err}}" placeholder it will be replaced with the original error message
func WrapError(e Error, format string, values ...interface{}) Error {
	msg := fmt.Sprintf(format, values...)
	if msg == "" {
		return e
	}
	return &err{
		cause: e.Cause(),
		pos:   e.Pos(),
		msg:   strings.Replace(msg, "{{err}}", e.Error(), -1),
	}
}
