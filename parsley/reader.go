// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Reader is a reader interface for parsing
//counterfeiter:generate . Reader
type Reader interface {
	Pos(int) Pos
	Remaining(Pos) int
	IsEOF(Pos) bool
}
