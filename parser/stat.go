// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

// Stat is the global statistics object
var Stat = Statistics{}

// Statistics contains parsing statistics
type Statistics struct {
	sumCallCount int
}

// RegisterCall registers a call
func (s *Statistics) RegisterCall() {
	s.sumCallCount++
}

// GetSumCallCount returns with the sum call count
func (s *Statistics) GetSumCallCount() int {
	return s.sumCallCount
}

// Reset resets the statistic counters
func (s *Statistics) Reset() {
	s.sumCallCount = 0
}
