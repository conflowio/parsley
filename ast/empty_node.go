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

// Type returns with an empty type
func (e EmptyNode) Type() string {
	return ""
}

// Value returns with nil
func (e EmptyNode) Value(ctx interface{}) (interface{}, parsley.Error) {
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
