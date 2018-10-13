package parsley

// File is an interface to translate a byte offset in a file to a position object
//go:generate counterfeiter . File
type File interface {
	Position(int) Position
	Pos(int) Pos
	Len() int
	SetOffset(int)
}
