// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package reader defines interfaces for an input reader and reader position
package reader

import (
	"fmt"
)

// Position is a given position in the reader
type Position interface {
	Pos() int
	fmt.Stringer
}

// Reader is a reader interface for parsing
type Reader interface {
	Clone() Reader
	Remaining() int
	Cursor() Position
	IsEOF() bool
	fmt.Stringer
}

// Error is an interface for handling errors with a position
type Error interface {
	error
	Cause() error
	Msg() string
	Pos() Position
}
