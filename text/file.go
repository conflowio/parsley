// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/conflowio/parsley/parsley"
)

// File contains the contents of a file and the line offsets for quick line+column lookup
type File struct {
	filename string
	data     []byte
	lines    []int
	len      int
	offset   int
}

// NewFile creates a new file object
func NewFile(filename string, data []byte) *File {
	f := &File{
		filename: filename,
		data:     bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1),
		offset:   1,
	}
	f.len = len(f.data)
	return f
}

// ReadFile reads a file and creates a File object
func ReadFile(filename string) (*File, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can not read %s", filename)
	}
	return NewFile(filename, data), nil
}

// SetLinesForContent sets the line offsets
func (f *File) setLines() {
	f.lines = []int{0}
	for offset, b := range f.data {
		if b == '\n' {
			f.lines = append(f.lines, offset+1)
		}
	}
}

// Len returns with the length of the file in bytes
func (f *File) Len() int {
	return f.len
}

// SetOffset set the offset of this file related to a file set
func (f *File) SetOffset(offset int) {
	f.offset = offset
}

// Position returns with a Position object for the given offset
func (f *File) Position(pos int) parsley.Position {
	if pos > f.len {
		return parsley.NilPosition
	}
	if f.lines == nil {
		f.setLines()
	}
	i := sort.Search(len(f.lines), func(i int) bool { return f.lines[i] > pos }) - 1
	return &Position{
		Filename: f.filename,
		Line:     i + 1,
		Column:   pos - f.lines[i] + 1,
	}
}

// Pos returns with a global offset in a file set
func (f *File) Pos(pos int) parsley.Pos {
	return parsley.Pos(f.offset + pos)
}
