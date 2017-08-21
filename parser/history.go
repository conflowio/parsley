// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

type storedResult struct {
	curtailingParsers data.IntSet
	resultSet         ResultSet
	err               reader.Error
	leftRecCtx        data.IntMap
}

// History records information about parser calls
type History struct {
	callCount int
	results   map[int]map[int]storedResult
}

// NewHistory creates a history instance
func NewHistory() *History {
	return &History{
		results: make(map[int]map[int]storedResult),
	}
}

// RegisterResults registers a parser result for a certain position
func (h *History) RegisterResults(parserIndex int, pos int, curtailingParsers data.IntSet, resultSet ResultSet, err reader.Error, leftRecCtx data.IntMap) {
	if _, ok := h.results[parserIndex]; !ok {
		h.results[parserIndex] = make(map[int]storedResult)
	}
	h.results[parserIndex][pos] = storedResult{curtailingParsers, resultSet, err, leftRecCtx}
}

// GetResults return with a previously saved result
func (h *History) GetResults(parserIndex int, pos int, leftRecCtx data.IntMap) (data.IntSet, ResultSet, reader.Error, bool) {
	storedResult, found := h.results[parserIndex][pos]
	if !found {
		return data.EmptyIntSet(), nil, nil, false
	}

	for key := range storedResult.leftRecCtx.Keys() {
		if storedResult.leftRecCtx.Get(key) > leftRecCtx.Get(key) {
			return data.EmptyIntSet(), nil, nil, false
		}
	}

	return storedResult.curtailingParsers, storedResult.resultSet, storedResult.err, true
}

// RegisterCall registers a call
func (h *History) RegisterCall() {
	h.callCount++
}

// CallCount returns with the call count
func (h *History) CallCount() int {
	return h.callCount
}
