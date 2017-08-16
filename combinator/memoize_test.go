// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/mock"
)

var memoizeFunc = parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
	return data.EmptyIntSet(), nil, nil
})

type historyMock struct {
	mock.Mock
}

func (h *historyMock) Memoize(p parser.Parser) parser.Func {
	h.Called(p)
	return memoizeFunc
}

type parserMock struct {
}

func (p parserMock) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
	return data.EmptyIntSet(), nil, nil
}

func TestMemoizeShouldCallHistoryMethod(t *testing.T) {
	h := new(historyMock)
	p := new(parserMock)
	h.On("Memoize", p).Return(memoizeFunc)
	combinator.Memoize("x", h, p)
	h.AssertExpectations(t)
}
