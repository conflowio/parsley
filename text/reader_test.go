// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/opsidian/parsley/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Let's read from a byte array with a regular expression
func ExampleReader() {
	r := text.NewReader([]byte("abcd"), true)
	matches, _, _ := r.ReadMatch("ab|cd", false)
	fmt.Println(matches[0])
	// Output: ab
}

func TestNewPosition(t *testing.T) {
	p := text.NewPosition(1, 2, 3)
	assert.Equal(t, "", p.Filename())
	assert.Equal(t, 1, p.Pos())
	assert.Equal(t, 2, p.Line())
	assert.Equal(t, 3, p.Col())
	assert.Equal(t, "2:3", p.String())
}

func TestNewFilePosition(t *testing.T) {
	p := text.NewFilePosition("file.name", 1, 2, 3)
	assert.Equal(t, "file.name", p.Filename())
	assert.Equal(t, 1, p.Pos())
	assert.Equal(t, 2, p.Line())
	assert.Equal(t, 3, p.Col())
	assert.Equal(t, "file.name:2:3", p.String())
}

func TestEmptyReader(t *testing.T) {
	r := text.NewReader([]byte{}, true)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
	assert.Equal(t, 0, r.Remaining())
	assert.True(t, r.IsEOF())
	_, _, err := r.ReadRune()
	assert.Exactly(t, io.EOF, err)
}

func TestNewReaderNotIgnoringWhitespacesShouldKeepWhitespaces(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\t foo\r\n\t "), false)
	assert.Equal(t, 12, r.Remaining())
	ch, _, _ := r.ReadRune()
	assert.Equal(t, ' ', ch)
}

// This test was introduced as the reader was originally trimming the starting whitespaces
func TestNewReaderShouldNotTrimInput(t *testing.T) {
	r := text.NewReader([]byte(" foo"), true)
	assert.Equal(t, 4, r.Remaining())
	ch, _, _ := r.ReadRune()
	assert.Equal(t, ' ', ch)
}

func TestCloneShouldCreateReaderWithSameParams(t *testing.T) {
	r := text.NewReader([]byte("ab\ncd\nef"), true)
	r.ReadMatch("ab\nc", false)
	rc := r.Clone().(*text.Reader)

	assert.Equal(t, r.Remaining(), rc.Remaining())
	assert.Equal(t, r.Cursor(), rc.Cursor())
	assert.Equal(t, r.IsEOF(), rc.IsEOF())

	rc.ReadMatch("d\nef", false)

	assert.Equal(t, 4, r.Remaining())
	assert.Equal(t, 0, rc.Remaining())
	assert.Equal(t, text.NewPosition(4, 2, 2), r.Cursor())
	assert.Equal(t, text.NewPosition(8, 3, 3), rc.Cursor())
	assert.False(t, r.IsEOF())
	assert.True(t, rc.IsEOF())
}

func TestReadRuneShouldReturnWithASCIICharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), true)
	ch, size, err := r.ReadRune()
	assert.Equal(t, 'a', ch)
	assert.Equal(t, 1, size)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(1, 1, 2), r.Cursor())
}

func TestReadRuneShouldReturnWithUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	ch, size, err := r.ReadRune()
	assert.Equal(t, '🍕', ch)
	assert.Equal(t, 4, size)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(4, 1, 2), r.Cursor())
}

func TestReadRuneShouldReturnErrorForInvalidUtfCharacter(t *testing.T) {
	r := text.NewReader([]byte("\xc3\x28"), true)
	_, _, err := r.ReadRune()
	assert.Error(t, err)
}

func TestReadRuneShouldReturnErrorIfNoMoreCharsLeft(t *testing.T) {
	var err error
	r := text.NewReader([]byte("a"), true)
	_, _, err = r.ReadRune()
	assert.Nil(t, err)
	_, _, err = r.ReadRune()
	assert.Exactly(t, io.EOF, err)
}

func TestReadRuneShouldFollowLinesAndColumns(t *testing.T) {
	r := text.NewReader([]byte("a\nb"), true)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())

	r.ReadRune()
	assert.Equal(t, text.NewPosition(1, 1, 2), r.Cursor())

	r.ReadRune()
	assert.Equal(t, text.NewPosition(2, 2, 1), r.Cursor())

	r.ReadRune()
	assert.Equal(t, text.NewPosition(3, 2, 2), r.Cursor())
}

func TestPeekRuneShouldReturnWithASCIICharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), true)
	ch, size, err := r.PeekRune()
	assert.Equal(t, 'a', ch)
	assert.Equal(t, 1, size)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestPeekRuneShouldReturnWithUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	ch, size, err := r.PeekRune()
	assert.Equal(t, '🍕', ch)
	assert.Equal(t, 4, size)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestPeekRuneShouldReturnErrorIfNoMoreCharsLeft(t *testing.T) {
	var err error
	r := text.NewReader([]byte(""), true)
	_, _, err = r.PeekRune()
	assert.Exactly(t, io.EOF, err)
}

func TestPeekRuneShouldReturnErrorForInvalidUtfCharacter(t *testing.T) {
	r := text.NewReader([]byte("\xc3\x28"), true)
	_, _, err := r.PeekRune()
	assert.Error(t, err)
}

func TestReadMatchShouldAlwaysMatchTheBeginning(t *testing.T) {
	r := text.NewReader([]byte("abc"), true)
	matches, _, ok := r.ReadMatch("x", false)
	assert.False(t, ok)
	assert.Nil(t, matches)
}

func TestReadMatchShouldAllPartsOfCompositeFromTheBeginning(t *testing.T) {
	r := text.NewReader([]byte("abcd"), true)
	matches, _, ok := r.ReadMatch("ab|cd", false)
	require.True(t, ok)
	assert.Equal(t, "ab", matches[0])

	r = text.NewReader([]byte("abcd"), true)
	matches, _, ok = r.ReadMatch("xx|cd", false)
	assert.False(t, ok)
	assert.Nil(t, matches)
}

func TestReadMatchShouldReturnMatchAndSubmatches(t *testing.T) {
	r := text.NewReader([]byte("123abcDEF"), true)
	matches, pos, ok := r.ReadMatch("(\\d+)([a-z]+)([A-Z]+)", false)
	require.True(t, ok)
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "123abcDEF", matches[0])
	assert.Equal(t, "123", matches[1])
	assert.Equal(t, "abc", matches[2])
	assert.Equal(t, "DEF", matches[3])
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
}

func TestReadMatchShouldReturnOnlyMainMatchIfNoCatchGroups(t *testing.T) {
	r := text.NewReader([]byte("abc"), true)
	matches, _, ok := r.ReadMatch("\\w+", false)
	require.True(t, ok)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
}

func TestReadMatchShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\tabc"), true)
	matches, pos, ok := r.ReadMatch("[a-z]+", false)
	require.True(t, ok)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
	assert.Equal(t, text.NewPosition(7, 2, 5), r.Cursor())
	assert.Equal(t, text.NewPosition(4, 2, 2), pos)
}

func TestReadMatchShouldNotIgnoreWhitespacesIfNotSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\tabc"), false)
	matches, _, ok := r.ReadMatch("[a-z]+", false)
	assert.False(t, ok)
	assert.Nil(t, matches)

	matches2, pos, ok := r.ReadMatch("\\s+[a-z]+", false)
	require.True(t, ok)
	assert.Equal(t, 1, len(matches2))
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
}

func TestReadMatchShouldIncludeWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\tabc"), true)
	matches, pos, ok := r.ReadMatch("\\s+[a-z]+", true)
	require.True(t, ok)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
}

func TestReadMatchShouldReturnFalseIfNoMatch(t *testing.T) {
	r := text.NewReader([]byte(" 123"), true)
	matches, pos, ok := r.ReadMatch("[a-z]+", false)
	assert.False(t, ok)
	assert.Nil(t, pos)
	assert.Nil(t, matches)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestReadMatchShouldFollowLinesAndColumns(t *testing.T) {
	r := text.NewReader([]byte("a\nb"), false)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())

	r.ReadMatch("(?s).", false)
	assert.Equal(t, text.NewPosition(1, 1, 2), r.Cursor())

	r.ReadMatch("(?s).", false)
	assert.Equal(t, text.NewPosition(2, 2, 1), r.Cursor())

	r.ReadMatch("(?s).", false)
	assert.Equal(t, text.NewPosition(3, 2, 2), r.Cursor())
}

func TestReadMatchShouldHandleUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	matches, pos, ok := r.ReadMatch(".*", false)
	require.True(t, ok)
	assert.Equal(t, []string{"🍕"}, matches)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(4, 1, 2), r.Cursor())
}

func TestPeekMatchShouldMatchButNotMoveCursor(t *testing.T) {
	r := text.NewReader([]byte("abc"), true)
	expectedPos := r.Cursor()
	matches, ok := r.PeekMatch("\\w+")
	require.True(t, ok)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
	assert.Equal(t, expectedPos, r.Cursor())
}

func TestPeekMatchShouldReturnMatchAndSubmatches(t *testing.T) {
	r := text.NewReader([]byte("123abcDEF"), true)
	matches, ok := r.PeekMatch("(\\d+)([a-z]+)([A-Z]+)")
	require.True(t, ok)
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "123abcDEF", matches[0])
	assert.Equal(t, "123", matches[1])
	assert.Equal(t, "abc", matches[2])
	assert.Equal(t, "DEF", matches[3])
}

func TestPeekMatchShouldReturnNilIfNoMatch(t *testing.T) {
	r := text.NewReader([]byte("123"), true)
	matches, ok := r.PeekMatch("[a-z]+")
	assert.False(t, ok)
	assert.Nil(t, matches)
}

func TestPeekMatchShouldNotIgnoreWhitespacesEvenIfSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\tabc"), true)
	matches, ok := r.PeekMatch("[a-z]+")
	assert.False(t, ok)
	assert.Nil(t, matches)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestStringShouldReturnNonEmptyString(t *testing.T) {
	r := text.NewReader([]byte("ab"), true)
	assert.NotEmpty(t, r.String())
}

func TestReadfShouldReturnResultAndPos(t *testing.T) {
	r := text.NewReader([]byte("123abcDEF"), true)
	reader := func(b []byte) (string, int, bool) {
		assert.Equal(t, []byte("123abcDEF"), b)
		return "NEXT: " + string(b[:3]), 3, true
	}

	result, pos, ok := r.Readf(reader, false)
	require.True(t, ok)
	assert.Equal(t, "NEXT: 123", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
}

func TestReadfShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n123abcd"), true)
	reader := func(b []byte) (string, int, bool) {
		assert.Equal(t, []byte("123abcd"), b)
		return "NEXT: " + string(b[:3]), 3, true
	}
	result, pos, ok := r.Readf(reader, false)
	require.True(t, ok)
	assert.Equal(t, "NEXT: 123", result)
	assert.Equal(t, text.NewPosition(3, 2, 1), pos)
	assert.Equal(t, text.NewPosition(6, 2, 4), r.Cursor())
}

func TestReadfShouldNotIgnoreWhitespacesIfNotSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n123abcd"), false)
	reader := func(b []byte) (string, int, bool) {
		assert.Equal(t, []byte(" \r\n123abcd"), b)
		return "NEXT: " + string(b[:3]), 3, true
	}
	result, pos, ok := r.Readf(reader, false)
	require.True(t, ok)
	assert.Equal(t, "NEXT:  \r\n", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(3, 2, 1), r.Cursor())
}

func TestReadfShouldIncludeWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n123abcd"), true)
	reader := func(b []byte) (string, int, bool) {
		assert.Equal(t, []byte(" \r\n123abcd"), b)
		return "NEXT: " + string(b[:3]), 3, true
	}
	result, pos, ok := r.Readf(reader, true)
	require.True(t, ok)
	assert.Equal(t, "NEXT:  \r\n", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(3, 2, 1), r.Cursor())
}

func TestReadfShouldReturnFalseIfNoMatch(t *testing.T) {
	r := text.NewReader([]byte("123"), true)
	reader := func(b []byte) (string, int, bool) {
		return "", 0, false
	}
	result, pos, ok := r.Readf(reader, false)
	assert.False(t, ok)
	assert.Equal(t, "", result)
	assert.Nil(t, pos)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestReadfShouldFollowLinesAndColumns(t *testing.T) {
	r := text.NewReader([]byte("a\nb"), false)
	reader := func(b []byte) (string, int, bool) {
		return "NEXT: " + string(b[:1]), 1, true
	}

	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())

	r.Readf(reader, false)
	assert.Equal(t, text.NewPosition(1, 1, 2), r.Cursor())

	r.Readf(reader, false)
	assert.Equal(t, text.NewPosition(2, 2, 1), r.Cursor())

	r.Readf(reader, false)
	assert.Equal(t, text.NewPosition(3, 2, 2), r.Cursor())
}

func TestReadfShouldHandleUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	reader := func(b []byte) (string, int, bool) {
		r, size := utf8.DecodeRuneInString(string(b))
		return string(r), size, true
	}
	result, pos, ok := r.Readf(reader, false)
	require.True(t, ok)
	assert.Equal(t, "🍕", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(4, 1, 2), r.Cursor())
}

func TestIsEOFShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte(" "), true)
	assert.True(t, r.IsEOF())
}

func TestIsEOFShouldReturnFalseIfNotAtTheEnd(t *testing.T) {
	r := text.NewReader([]byte(" "), false)
	assert.False(t, r.IsEOF())
}

func TestNewFileReader(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "text_reader_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	r, err := text.NewFileReader(tmpfile.Name(), true)
	assert.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, text.NewFilePosition(tmpfile.Name(), 0, 1, 1), r.Cursor())
	assert.Equal(t, 0, r.Remaining())
	assert.True(t, r.IsEOF())
}

func TestNewFileReaderWithNonexistingFile(t *testing.T) {
	r, err := text.NewFileReader("non-existing.file", true)
	assert.Error(t, err)
	assert.Nil(t, r)
}
