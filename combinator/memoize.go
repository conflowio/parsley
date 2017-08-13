// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/parser"
)

type memoizer interface {
	Memoize(p parser.Parser) parser.Func
}

// Memoize handles result cache and curtailing left recursion
//
// Deprecated: please use the Memoize method on the history object
func Memoize(name string, h memoizer, p parser.Parser) parser.Func {
	return h.Memoize(p)
}
