// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Result is a stored parser result
type Result struct {
	LeftRecCtx        data.IntMap
	CurtailingParsers data.IntSet
	Error             Error
	Node              Node
}

// ResultCache records information about parser calls
type ResultCache map[int]map[Pos]*Result

// NewResultCache creates a history instance
func NewResultCache() ResultCache {
	return make(map[int]map[Pos]*Result)
}

// Save registers a parser result for a certain position
func (rc ResultCache) Save(parserIndex int, pos Pos, result *Result) {
	if _, ok := rc[parserIndex]; !ok {
		rc[parserIndex] = make(map[Pos]*Result)
	}
	rc[parserIndex][pos] = result
}

// Get return with a previously saved result
func (rc ResultCache) Get(parserIndex int, pos Pos, leftRecCtx data.IntMap) (*Result, bool) {
	result, found := rc[parserIndex][pos]
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
