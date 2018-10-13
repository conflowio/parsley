// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/parsley"
)

// SepBy applies the given value parser zero or more times separated by the separator parser
func SepBy(name string, valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(name, valueP, sepP, true)
}

// SepBy1 applies the given value parser one or more times separated by the separator parser
func SepBy1(name string, valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(name, valueP, sepP, false)
}

func newSepBy(name string, valueP parsley.Parser, sepP parsley.Parser, allowEmpty bool) *Recursive {
	lookup := func(i int) parsley.Parser {
		if i%2 == 0 {
			return valueP
		} else {
			return sepP
		}
	}
	lenCheck := func(len int) bool {
		return (len == 0 && allowEmpty) || len%2 == 1
	}
	return NewRecursive("SEP_BY", name, lookup, lenCheck)
}
