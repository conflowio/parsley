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
type Error interface {
	error
	Pos() Pos
}

type err struct {
	msg string
	pos Pos
}

// NewError creates a new error with the given position
// If the passed error is already a parsley.Error it returns the original error
// as it should have already the correct position.
func NewError(pos Pos, format string, values ...interface{}) Error {
	return &err{
		msg: fmt.Sprintf(format, values...),
		pos: pos,
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

// WrapError wraps the given error in a error
// If format contains the "{{err}}" placeholder it will be replaced with the original error message
func WrapError(e Error, format string, values ...interface{}) Error {
	msg := fmt.Sprintf(format, values...)
	if msg == "" {
		return e
	}
	return &err{
		msg: strings.Replace(msg, "{{err}}", e.Error(), -1),
		pos: e.Pos(),
	}
}
