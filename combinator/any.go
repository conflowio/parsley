// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Any tries all the given parsers independently and merges the results
func Any(desc string, parsers ...parsley.Parser) parsley.ParserFunc {
	if parsers == nil {
		panic("No parsers were given")
	}
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		cp := data.EmptyIntSet()
		var nodes []parsley.Node
		var err parsley.Error
		for _, p := range parsers {
			h.RegisterCall()
			cp2, nodes2, err2 := p.Parse(h, leftRecCtx, r, pos)
			cp = cp.Union(cp2)
			nodes = append(nodes, nodes2...)
			if err2 != nil && (err == nil || err2.Pos() >= err.Pos()) {
				err = err2
			}
		}
		if desc != "" && err != nil && err.Pos() == r.Pos(pos) {
			err = parsley.NewError(r.Pos(pos), "was expecting %s", desc)
		}
		return cp, nodes, err
	})
}
