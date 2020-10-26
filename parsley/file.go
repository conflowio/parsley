// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// File is an interface to translate a byte offset in a file to a position object
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . File
type File interface {
	Position(int) Position
	Pos(int) Pos
	Len() int
	SetOffset(int)
}
