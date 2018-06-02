// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"fmt"

	"github.com/jinzhu/inflection"
	"github.com/opsidian/parsley/parsley"
)

// SepBy applies the given value parser zero or more times separated by the separator parser
func SepBy(valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(valueP, sepP, true, false)
}

// SepByOrValue applies the given value parser zero or more times separated by the separator parser
// If there is only one value then the value node will be returned
func SepByOrValue(valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(valueP, sepP, true, true)
}

// SepBy1 applies the given value parser one or more times separated by the separator parser
func SepBy1(valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(valueP, sepP, false, false)
}

// SepByOrValue1 applies the given value parser one or more times separated by the separator parser
// If there is only one value then the value node will be returned
func SepByOrValue1(valueP parsley.Parser, sepP parsley.Parser) *Recursive {
	return newSepBy(valueP, sepP, false, true)
}

func newSepBy(valueP parsley.Parser, sepP parsley.Parser, allowEmpty bool, returnSingle bool) *Recursive {
	name := func() string {
		return fmt.Sprintf("%s separated by %s", inflection.Plural(valueP.Name()), inflection.Plural(sepP.Name()))
	}
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
	return NewRecursive("SEP_BY", name, returnSingle, lookup, lenCheck)
}
