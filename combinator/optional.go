// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Optional returns the parser's matches or an empty match if it fails
func Optional(p parsley.Parser) *parser.NamedFunc {
	return Choice("", p, parser.Empty())
}
