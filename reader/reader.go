// Package reader defines interfaces for an input reader and reader position
package reader

import "fmt"

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
