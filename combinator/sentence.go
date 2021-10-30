// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/conflowio/parsley/ast/interpreter"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
)

// Sentence matches the given parser until the end of input
func Sentence(p parsley.Parser) *Sequence {
	return SeqOf(p, parser.End()).Bind(interpreter.Select(0))
}
