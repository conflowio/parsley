// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	pos := new(mocks.Position)
	pos.On("String").Return("POS")
	err := parser.NewError(pos, "ERR")
	assert.Equal(t, pos, err.Pos())
	assert.Equal(t, "ERR", err.Msg())
	assert.Equal(t, "ERR at POS", err.Error())
}

func TestWrapErrorShouldWrapErr(t *testing.T) {
	cause := errors.New("ERR")
	pos := new(mocks.Position)
	pos.On("String").Return("POS")
	err := parser.WrapError(pos, cause)
	assert.Equal(t, pos, err.Pos())
	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, "ERR", err.Msg())
	assert.Equal(t, "ERR at POS", err.Error())
}

func TestWrapErrorShouldNotDoubleWrap(t *testing.T) {
	cause := parser.NewError(new(mocks.Position), "ERR")
	err := parser.WrapError(new(mocks.Position), cause)
	assert.Equal(t, cause, err)
}
