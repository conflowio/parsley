// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
