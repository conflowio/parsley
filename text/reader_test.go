package text_test

import (
	"io"
	"testing"
	"unicode/utf8"

	"github.com/opsidian/parsley/text"
	"github.com/stretchr/testify/assert"
)

func assertCursor(t *testing.T, pos int, line int, col int, r *text.Reader) {
	assert.Equal(t, text.NewPosition(pos, line, col), r.Cursor())
}

func TestPositionMethods(t *testing.T) {
	p := text.NewPosition(1, 2, 3)
	assert.Equal(t, 1, p.Pos())
	assert.Equal(t, 2, p.Line())
	assert.Equal(t, 3, p.Col())
	assert.NotEmpty(t, p.String())
}

func TestEmptyReader(t *testing.T) {
	r := text.NewReader([]byte{}, true)
	assertCursor(t, 0, 1, 1, r)
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

func TestNewReaderIgnoringWhitespacesShouldTrimText(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\t foo\r\n\t "), true)
	assert.Equal(t, 3, r.Remaining())
	ch, _, _ := r.ReadRune()
	assert.Equal(t, 'f', ch)
}

func TestCloneShouldCreateReaderWithSameParams(t *testing.T) {
	r := text.NewReader([]byte("ab\ncd\nef"), true)
	r.ReadMatch("ab\nc")
	rc := r.Clone().(*text.Reader)

	assert.Equal(t, r.Remaining(), rc.Remaining())
	assert.Equal(t, r.Cursor(), rc.Cursor())
	assert.Equal(t, r.IsEOF(), rc.IsEOF())

	rc.ReadMatch("d\nef")

	assert.Equal(t, 4, r.Remaining())
	assert.Equal(t, 0, rc.Remaining())
	assertCursor(t, 4, 2, 2, r)
	assertCursor(t, 8, 3, 3, rc)
	assert.False(t, r.IsEOF())
	assert.True(t, rc.IsEOF())

}

func TestReadRuneShouldReturnWithCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	ch, size, err := r.ReadRune()
	assert.Equal(t, '🍕', ch)
	assert.Equal(t, 4, size)
	assert.Nil(t, err)
	assertCursor(t, 4, 1, 2, r)
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
	assertCursor(t, 0, 1, 1, r)

	r.ReadRune()
	assertCursor(t, 1, 1, 2, r)

	r.ReadRune()
	assertCursor(t, 2, 2, 1, r)

	r.ReadRune()
	assertCursor(t, 3, 2, 2, r)
}

func TestReadMatchShouldAlwaysMatchTheBeginning(t *testing.T) {
	r := text.NewReader([]byte("abc"), true)
	matches, _ := r.ReadMatch("x")
	assert.Nil(t, matches)
}

func TestReadMatchShouldAllPartsOfCompositeFromTheBeginning(t *testing.T) {
	r := text.NewReader([]byte("abcd"), true)
	matches, _ := r.ReadMatch("ab|cd")
	assert.Equal(t, "ab", matches[0])

	r = text.NewReader([]byte("abcd"), true)
	matches, _ = r.ReadMatch("xx|cd")
	assert.Nil(t, matches)
}

func TestReadMatchShouldReturnMatchAndSubmatches(t *testing.T) {
	r := text.NewReader([]byte("123abcDEF"), true)
	matches, pos := r.ReadMatch("(\\d+)([a-z]+)([A-Z]+)")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "123abcDEF", matches[0])
	assert.Equal(t, "123", matches[1])
	assert.Equal(t, "abc", matches[2])
	assert.Equal(t, "DEF", matches[3])
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
}

func TestReadMatchShouldReturnOnlyMainMatchIfNoCatchGroups(t *testing.T) {
	r := text.NewReader([]byte("abc"), true)
	matches, _ := r.ReadMatch("\\w+")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
}

func TestReadMatchShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte("x \r\n\tabc"), true)
	r.ReadRune()
	matches, pos := r.ReadMatch("[a-z]+")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "abc", matches[0])
	assertCursor(t, 8, 2, 5, r)
	assert.Equal(t, text.NewPosition(5, 2, 2), pos)
}

func TestReadMatchShouldNotIgnoreWhitespacesIfNotSet(t *testing.T) {
	r := text.NewReader([]byte(" \r\n\tabc"), false)
	matches, _ := r.ReadMatch("[a-z]+")
	assert.Nil(t, matches)

	matches2, pos := r.ReadMatch("\\s+[a-z]+")
	assert.Equal(t, 1, len(matches2))
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
}

func TestReadMatchShouldReturnNilIfNoMatch(t *testing.T) {
	r := text.NewReader([]byte("123"), true)
	matches, _ := r.ReadMatch("[a-z]+")
	assert.Nil(t, matches)
}

func TestReadMatchShouldFollowLinesAndColumns(t *testing.T) {
	r := text.NewReader([]byte("a\nb"), false)
	assertCursor(t, 0, 1, 1, r)

	r.ReadMatch("(?s).")
	assertCursor(t, 1, 1, 2, r)

	r.ReadMatch("(?s).")
	assertCursor(t, 2, 2, 1, r)

	r.ReadMatch("(?s).")
	assertCursor(t, 3, 2, 2, r)
}

func TestReadMatchShouldHandleUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	matches, pos := r.ReadMatch(".*")
	assert.Equal(t, []string{"🍕"}, matches)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assertCursor(t, 4, 1, 2, r)
}

func TestStringShouldReturnNonEmptyString(t *testing.T) {
	r := text.NewReader([]byte("ab"), true)
	r.ReadRune()
	assert.NotEmpty(t, r.String())
}

func TestWithCursorShouldSetPositionAndCharCount(t *testing.T) {
	r := text.NewReader([]byte("abcd"), true)
	r2 := r.WithCursor(3, 1, 4)
	assert.Equal(t, text.NewPosition(3, 1, 4), r2.Cursor())
	assert.Equal(t, 1, r2.Remaining())
}

func TestReadfShouldReturnResultAndPos(t *testing.T) {
	r := text.NewReader([]byte("123abcDEF"), true)
	reader := func(b []byte) (string, int) {
		assert.Equal(t, []byte("123abcDEF"), b)
		return "NEXT: " + string(b[:3]), 3
	}

	result, pos := r.Readf(reader)
	assert.Equal(t, "NEXT: 123", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
}

func TestReadfShouldIgnoreWhitespacesIfSet(t *testing.T) {
	r := text.NewReader([]byte("x \r\n123abcd"), true)
	r.ReadRune()
	reader := func(b []byte) (string, int) {
		assert.Equal(t, []byte("123abcd"), b)
		return "NEXT: " + string(b[:3]), 3
	}
	result, pos := r.Readf(reader)
	assert.Equal(t, "NEXT: 123", result)
	assert.Equal(t, text.NewPosition(4, 2, 1), pos)
	assert.Equal(t, text.NewPosition(7, 2, 4), r.Cursor())
}

func TestReadfShouldNotIgnoreWhitespacesIfNotSet(t *testing.T) {
	r := text.NewReader([]byte("x \r\n123abcd"), false)
	r.ReadRune()
	reader := func(b []byte) (string, int) {
		assert.Equal(t, []byte(" \r\n123abcd"), b)
		return "NEXT: " + string(b[:3]), 3
	}
	result, pos := r.Readf(reader)
	assert.Equal(t, "NEXT:  \r\n", result)
	assert.Equal(t, text.NewPosition(1, 1, 2), pos)
	assert.Equal(t, text.NewPosition(4, 2, 1), r.Cursor())
}

func TestReadfShouldReturnEmptyStringIfNoMatch(t *testing.T) {
	r := text.NewReader([]byte("123"), true)
	reader := func(b []byte) (string, int) {
		return "", 0
	}
	result, pos := r.Readf(reader)
	assert.Equal(t, "", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestReadfShouldFollowLinesAndColumns(t *testing.T) {
	r := text.NewReader([]byte("a\nb"), false)
	reader := func(b []byte) (string, int) {
		return "NEXT: " + string(b[:1]), 1
	}

	assertCursor(t, 0, 1, 1, r)

	r.Readf(reader)
	assertCursor(t, 1, 1, 2, r)

	r.Readf(reader)
	assertCursor(t, 2, 2, 1, r)

	r.Readf(reader)
	assertCursor(t, 3, 2, 2, r)
}

func TestReadfShouldHandleUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	reader := func(b []byte) (string, int) {
		r, size := utf8.DecodeRuneInString(string(b))
		return string(r), size
	}
	result, pos := r.Readf(reader)
	assert.Equal(t, "🍕", result)
	assert.Equal(t, text.NewPosition(0, 1, 1), pos)
	assert.Equal(t, text.NewPosition(4, 1, 2), r.Cursor())
}
