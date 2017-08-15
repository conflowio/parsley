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
	pos   reader.Position
}

// NewError creates a new parser error instance
func NewError(pos reader.Position, format string, values ...interface{}) Error {
	return err{
		cause: fmt.Errorf(format, values...),
		pos:   pos,
	}
}

// WrapError wraps the given error in a parser error
// If the cause is already a parser error that will be returned as is
func WrapError(pos reader.Position, cause error) Error {
	if parserErr, ok := cause.(Error); ok {
		return parserErr
	}
	return err{
		cause: cause,
		pos:   pos,
	}
}

// Cause returns with the original error
func (e err) Cause() error {
	return e.cause
}

// Msg returns with the error message
func (e err) Msg() string {
	return e.cause.Error()
}

// Error returns with the full error message including the position
func (e err) Error() string {
	return fmt.Sprintf("%s at %s", e.cause, e.pos)
}

// Pos returns with the error's position
func (e err) Pos() reader.Position {
	return e.pos
}
