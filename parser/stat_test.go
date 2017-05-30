package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCallShouldIncreaseSumCallCount(t *testing.T) {
	parser.Stat.Reset()
	assert.Equal(t, parser.Stat.GetSumCallCount(), 0)
	parser.Stat.RegisterCall()
	assert.Equal(t, parser.Stat.GetSumCallCount(), 1)
}

func TestResetShouldResetCounter(t *testing.T) {
	parser.Stat.Reset()
	parser.Stat.RegisterCall()
	assert.Equal(t, parser.Stat.GetSumCallCount(), 1)
	parser.Stat.Reset()
	assert.Equal(t, parser.Stat.GetSumCallCount(), 0)
}
