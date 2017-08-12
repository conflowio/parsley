// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"sync/atomic"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

type storedResult struct {
	curtailingParsers data.IntSet
	resultSet         ResultSet
	err               Error
	leftRecCtx        data.IntMap
}

// History records information about parser calls
type History struct {
	parserIndex int32
	results     map[int]map[int]storedResult
}

// NewHistory creates a history instance
func NewHistory() *History {
	return &History{
		parserIndex: 0,
		results:     make(map[int]map[int]storedResult),
	}
}

// Reset deletes the collected data
func (h *History) Reset() {
	h.results = make(map[int]map[int]storedResult)
}

// RegisterResults registers a parser result for a certain position
func (h *History) RegisterResults(parserIndex int, pos int, curtailingParsers data.IntSet, resultSet ResultSet, err Error, leftRecCtx data.IntMap) {
	if _, ok := h.results[parserIndex]; !ok {
		h.results[parserIndex] = make(map[int]storedResult)
	}
	h.results[parserIndex][pos] = storedResult{curtailingParsers, resultSet, err, leftRecCtx}
}

// GetResults return with a previously saved result
func (h *History) GetResults(parserIndex int, pos int, leftRecCtx data.IntMap) (data.IntSet, ResultSet, Error, bool) {
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

// Memoize handles result cache and curtailing left recursion
func (h *History) Memoize(p Parser) Func {
	parserIndex := int(atomic.AddInt32(&h.parserIndex, 1))
	return Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet, Error) {
		cp, rs, err, found := h.GetResults(parserIndex, r.Cursor().Pos(), leftRecCtx)
		if found {
			return cp, rs, err
		}

		if leftRecCtx.Get(parserIndex) > r.Remaining()+1 {
			return data.NewIntSet(parserIndex), nil, nil
		}

		cp, rs, err = p.Parse(leftRecCtx.Inc(parserIndex), r)
		leftRecCtx = leftRecCtx.Filter(cp)

		h.RegisterResults(parserIndex, r.Cursor().Pos(), cp, rs, err, leftRecCtx)

		return cp, rs, err
	})
}
