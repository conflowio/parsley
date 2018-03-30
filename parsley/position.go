package parsley

import "fmt"

// Position is an interface to translate a file position to a string
//go:generate counterfeiter . Position
type Position interface {
	fmt.Stringer
}

type nilPosition int

func (np nilPosition) String() string {
	return "unknown"
}

// NilPosition represents an invalid position
const NilPosition = nilPosition(0)

// Pos is a global offset in a file set which can be translated into a concrete file position
type Pos int

// NilPos represents an invalid position
const NilPos = Pos(0)
