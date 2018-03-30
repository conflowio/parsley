package text

import (
	"fmt"
)

// Position is a text file position
type Position struct {
	Filename string
	Line     int
	Column   int
}

// NewPosition creates a new text position
func NewPosition(filename string, line int, column int) *Position {
	return &Position{
		Filename: filename,
		Line:     line,
		Column:   column,
	}
}

func (pos Position) String() string {
	if pos.Filename != "" {
		return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)
	}

	return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
}
