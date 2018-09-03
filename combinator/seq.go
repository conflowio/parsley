// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/parsley"
)

// Seq tries to apply all parsers after each other matching effectively a sequence of tokens
// and returns with all combinations of the results.
// Only matches are returned where all parsers were applied successfully.
func Seq(token string, name string, parsers ...parsley.Parser) *Recursive {
	return newSeq(token, name, false, len(parsers), parsers...)
}

// SeqTry tries to apply all parsers after each other matching effectively the longest possible
// sequences of tokens and returns with all combinations of the results.
// It needs to match the first parser at least
func SeqTry(token string, name string, parsers ...parsley.Parser) *Recursive {
	return newSeq(token, name, false, 1, parsers...)
}

// SeqTryOrValue tries to apply all parsers after each other matching effectively the longest possible
// sequences of tokens and returns with all combinations of the results.
// It needs to match the first parser at least.
// If only the first parser matches it will return the result of the first parser directly.
func SeqTryOrValue(token string, name string, parsers ...parsley.Parser) *Recursive {
	return newSeq(token, name, true, 1, parsers...)
}

func newSeq(token string, name string, returnSingle bool, min int, parsers ...parsley.Parser) *Recursive {
	namef := parsers[0].Name
	if name != "" {
		namef = func() string { return name }
	}
	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	l := len(parsers)
	lenCheck := func(len int) bool {
		return len >= min && len <= l
	}
	return NewRecursive(token, namef, returnSingle, lookup, lenCheck)
}
