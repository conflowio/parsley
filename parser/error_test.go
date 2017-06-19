package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := parser.NewError(test.NewPosition(1), "ERR1")
	assert.Equal(t, test.NewPosition(1), err.Pos())
	assert.Equal(t, "ERR1 at Pos{1}", err.Error())
}
