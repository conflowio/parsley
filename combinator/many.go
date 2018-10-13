// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/parsley"
)

// Many applies the  parser zero or more times
func Many(p parsley.Parser) *Recursive {
	return newMany(p, true)
}

// Many1 applies the parser one or more times
func Many1(p parsley.Parser) *Recursive {
	return newMany(p, false)
}

func newMany(p parsley.Parser, allowEmpty bool) *Recursive {
	lookup := func(i int) parsley.Parser {
		return p
	}
	lenCheck := func(len int) bool {
		return allowEmpty || len > 0
	}
	return NewRecursive("MANY", lookup, lenCheck)
}
