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

func TestNewErrorWithParams(t *testing.T) {
	pos := new(mocks.Position)
	pos.On("String").Return("POS")
	err := parser.NewError(pos, "ERR %s", "message")
	assert.Equal(t, "ERR message", err.Msg())
	assert.Equal(t, "ERR message at POS", err.Error())
}

func TestWrapErrorShouldWrapErr(t *testing.T) {
	cause := errors.New("WRAPPED ERR")
	pos := new(mocks.Position)
	pos.On("String").Return("POS")
	err := parser.WrapError(pos, cause, "ERR")
	assert.Equal(t, pos, err.Pos())
	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, "ERR", err.Msg())
	assert.Equal(t, "ERR at POS", err.Error())
}

func TestWrapErrorShouldWrapErrWithParams(t *testing.T) {
	cause := errors.New("WRAPPED ERR")
	pos := new(mocks.Position)
	pos.On("String").Return("POS")
	err := parser.WrapError(pos, cause, "ERR %s", "message")
	assert.Equal(t, "ERR message", err.Msg())
	assert.Equal(t, "ERR message at POS", err.Error())
}

func TestWrapErrorShouldUpdateErrorMessageForParserError(t *testing.T) {
	cause1 := errors.New("CAUSE 1")
	pos1 := new(mocks.Position)
	pos1.On("String").Return("POS 1")
	err1 := parser.WrapError(pos1, cause1, "ERR 1")

	pos2 := new(mocks.Position)
	pos2.On("String").Return("POS 2")

	err2 := parser.WrapError(pos2, err1, "ERR 2")
	assert.Equal(t, pos1, err2.Pos())
	assert.Equal(t, cause1, err2.Cause())
	assert.Equal(t, "ERR 2", err2.Msg())
	assert.Equal(t, "ERR 2 at POS 1", err2.Error())
}
