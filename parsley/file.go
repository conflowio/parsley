package parsley

import (
	"errors"
	"fmt"
	"sort"
)

// File is an interface to translate a byte offset in a file to a position object
//go:generate counterfeiter . File
type File interface {
	Position(int) Position
	Pos(int) Pos
	Len() int
	SetOffset(int)
}

// FileSet contains multiple files
type FileSet struct {
	pos    int
	files  []File
	offset []int
}

// NewFileSet creates a new file set
func NewFileSet(files ...File) *FileSet {
	fs := &FileSet{
		pos:    1,
		files:  []File{},
		offset: []int{},
	}
	for _, f := range files {
		fs.AddFile(f)
	}
	return fs
}

// AddFile adds a new file
func (fs *FileSet) AddFile(f File) {
	if f == nil {
		panic(errors.New("file can not be nil"))
	}

	f.SetOffset(fs.pos)
	fs.files = append(fs.files, f)
	fs.offset = append(fs.offset, fs.pos)
	fs.pos = fs.pos + f.Len() + 1
}

// Position returns with a position object for a given global position
func (fs *FileSet) Position(pos Pos) Position {
	if pos == 0 || int(pos) >= fs.pos {
		return NilPosition
	}

	i := sort.Search(len(fs.offset), func(i int) bool { return fs.offset[i] > int(pos) }) - 1
	return fs.files[i].Position(int(pos) - fs.offset[i])
}

// ErrorWithPosition creates an error with a human-readable position
func (fs *FileSet) ErrorWithPosition(err Error) error {
	pos := fs.Position(err.Pos())
	if pos == NilPosition {
		return err
	}
	return fmt.Errorf("%s at %s", err.Error(), pos.String())
}
