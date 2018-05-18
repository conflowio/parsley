// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// History records information about parser calls
type History struct {
	callCount int
	results   map[int]map[parsley.Pos]*parsley.Result
}

// NewHistory creates a history instance
func NewHistory() *History {
	return &History{
		results: make(map[int]map[parsley.Pos]*parsley.Result),
	}
}

// SaveResult registers a parser result for a certain position
func (h *History) SaveResult(parserIndex int, pos parsley.Pos, result *parsley.Result) {
	if _, ok := h.results[parserIndex]; !ok {
		h.results[parserIndex] = make(map[parsley.Pos]*parsley.Result)
	}
	h.results[parserIndex][pos] = result
}

// GetResult return with a previously saved result
func (h *History) GetResult(parserIndex int, pos parsley.Pos, leftRecCtx data.IntMap) (*parsley.Result, bool) {
	result, found := h.results[parserIndex][pos]
	if !found {
		return nil, false
	}

	for _, key := range result.LeftRecCtx.Keys() {
		if result.LeftRecCtx.Get(key) > leftRecCtx.Get(key) {
			return nil, false
		}
	}

	return result, true
}

// RegisterCall registers a call
func (h *History) RegisterCall() {
	h.callCount++
}

// CallCount returns with the call count
func (h *History) CallCount() int {
	return h.callCount
}
