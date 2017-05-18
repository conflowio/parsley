package reader_test

import (
	"io"
	"testing"

	"github.com/opsidian/parsec/reader"
	"github.com/stretchr/testify/assert"
)

func TestEmptyReader(t *testing.T) {
	r := reader.New([]byte{}, true)
	assert.Equal(t, 0, r.Position())
	assert.Equal(t, 0, r.CharsRemaining())
	assert.Equal(t, 1, r.Line())
	assert.Equal(t, 1, r.Column())
	assert.True(t, r.IsEOF())
	_, _, err := r.ReadRune()
	assert.Exactly(t, io.EOF, err)
}

func TestNewReaderNotIgnoringWhitespacesShouldKeepWhitespaces(t *testing.T) {
	r := reader.New([]byte(" \r\n\t foo\r\n\t "), false)
	assert.Equal(t, 12, r.CharsRemaining())
	ch, _, _ := r.ReadRune()
	assert.Equal(t, ' ', ch)
}

func TestNewReaderIgnoringWhitespacesShouldTrimText(t *testing.T) {
	r := reader.New([]byte(" \r\n\t foo\r\n\t "), true)
	assert.Equal(t, 3, r.CharsRemaining())
	ch, _, _ := r.ReadRune()
	assert.Equal(t, 'f', ch)
}

func TestCloneShouldCreateReaderWithSameParams(t *testing.T) {
	r := reader.New([]byte("ab\ncd\nef"), true)
	r.ReadMatch("^ab\nc")
	rc := r.Clone()

	assert.Equal(t, r.CharsRemaining(), rc.CharsRemaining())
	assert.Equal(t, r.Position(), rc.Position())
	assert.Equal(t, r.IsEOF(), rc.IsEOF())

	rc.ReadMatch("^d\nef")

	assert.Equal(t, 4, r.CharsRemaining())
	assert.Equal(t, 0, rc.CharsRemaining())
	assert.False(t, r.IsEOF())
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 2, r.Column())

	assert.Equal(t, 4, r.Position())
	assert.Equal(t, 8, rc.Position())
	assert.Equal(t, 3, rc.Line())
	assert.Equal(t, 3, rc.Column())
	assert.True(t, rc.IsEOF())

}

func TestReadRuneShouldReturnWithCharacter(t *testing.T) {
	r := reader.New([]byte("🍕"), true)
	ch, size, err := r.ReadRune()
	assert.Equal(t, '🍕', ch)
	assert.Equal(t, 4, size)
	assert.Nil(t, err)
}

func TestReadRuneShouldReturnErrorForInvalidUtfCharacter(t *testing.T) {
	r := reader.New([]byte("\xc3\x28"), true)
	_, _, err := r.ReadRune()
	assert.Error(t, err)
}

func TestReadRuneShouldReturnErrorIfNoMoreCharsLeft(t *testing.T) {
	var err error
	r := reader.New([]byte("a"), true)
	_, _, err = r.ReadRune()
	assert.Nil(t, err)
	_, _, err = r.ReadRune()
	assert.Exactly(t, io.EOF, err)
}

func TestReadRuneShouldFollowLinesAndColumns(t *testing.T) {
	r := reader.New([]byte("a\nb"), true)
	assert.Equal(t, 1, r.Line())
	assert.Equal(t, 1, r.Column())

	r.ReadRune()
	assert.Equal(t, 1, r.Line())
	assert.Equal(t, 2, r.Column())

	r.ReadRune()
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 1, r.Column())

	r.ReadRune()
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 2, r.Column())
}

func TestReadMatchShouldReturnErrorIfNotMatchingTheBeginning(t *testing.T) {
	r := reader.New([]byte("abc"), true)
	assert.Panics(t, func() { r.ReadMatch("x") })
}

func TestReadMatchShouldReturnMatchAndSubmatches(t *testing.T) {
	r := reader.New([]byte("123abcDEF"), true)
	matches, pos := r.ReadMatch("^(\\d+)([a-z]+)([A-Z]+)")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "123abcDEF", matches[0])
	assert.Equal(t, "123", matches[1])
	assert.Equal(t, "abc", matches[2])
	assert.Equal(t, "DEF", matches[3])
	assert.Equal(t, 0, pos)
}

func TestReadMatchShouldReturnOnlyMainMatchIfNoCatchGroups(t *testing.T) {
	r := reader.New([]byte("abc"), true)
	matches, _ := r.ReadMatch("^\\w+")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
}

func TestReadMatchShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := reader.New([]byte("x \r\n\tabc"), true)
	r.ReadRune()
	matches, pos := r.ReadMatch("^[a-z]+")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
	assert.Equal(t, 8, r.Position())
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 5, r.Column())
	assert.Equal(t, 5, pos)
}

func TestReadMatchShouldNotIgnoreWhitespacesIfNotSet(t *testing.T) {
	r := reader.New([]byte(" \r\n\tabc"), false)
	matches, _ := r.ReadMatch("^[a-z]+")
	assert.Nil(t, matches)

	matches2, pos := r.ReadMatch("^\\s+[a-z]+")
	assert.Equal(t, 1, len(matches2))
	assert.Equal(t, 0, pos)
}

func TestReadMatchShouldReturnNilIfNoMatch(t *testing.T) {
	r := reader.New([]byte("123"), true)
	matches, _ := r.ReadMatch("^[a-z]+")
	assert.Nil(t, matches)
}

func TestReadMatchShouldFollowLinesAndColumns(t *testing.T) {
	r := reader.New([]byte("a\nb"), false)
	assert.Equal(t, 1, r.Line())
	assert.Equal(t, 1, r.Column())

	r.ReadMatch("^(?s).")
	assert.Equal(t, 1, r.Line())
	assert.Equal(t, 2, r.Column())

	r.ReadMatch("^(?s).")
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 1, r.Column())

	r.ReadMatch("^(?s).")
	assert.Equal(t, 2, r.Line())
	assert.Equal(t, 2, r.Column())
}

func TestStringShouldReturnStatusString(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	r.ReadRune()
	assert.Equal(t, "Reader{pos: 1, 1 chars left}\n", r.String())
}
