// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"
)

// EmptyNode represents an empty node
type EmptyNode parsley.Pos

// Token returns with EMPTY
func (e EmptyNode) Token() string {
	return "EMPTY"
}

// Schema returns nil
func (e EmptyNode) Schema() interface{} {
	return nil
}

// Value returns with nil
func (e EmptyNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	return nil, nil
}

// Pos returns with the position of the node
func (e EmptyNode) Pos() parsley.Pos {
	return parsley.Pos(e)
}

// ReaderPos returns the reader position
func (e EmptyNode) ReaderPos() parsley.Pos {
	return parsley.Pos(e)
}

// String returns with a string representation of the node
func (e EmptyNode) String() string {
	return fmt.Sprintf("%s{%d}", e.Token(), e.Pos())
}
