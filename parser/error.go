package parser

import (
	"fmt"

	"github.com/opsidian/parsley/reader"
)

// Error is a parser error interface
type Error interface {
	error
	Pos() reader.Position
}

// err is a parser error
type err struct {
	msg string
	pos reader.Position
}

// NewError creates a new parser error instance
func NewError(pos reader.Position, format string, values ...interface{}) Error {
	return err{
		msg: fmt.Sprintf(format, values...),
		pos: pos,
	}
}

// Error returns with the error message
func (e err) Error() string {
	return fmt.Sprintf("%s at %s", e.msg, e.pos)
}

// Pos returns with the error's position
func (e err) Pos() reader.Position {
	return e.pos
}
