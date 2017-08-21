// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Choice tries to apply the given parsers until one of them succeeds
func Choice(desc string, parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		cur := r.Cursor()
		cp := parser.NoCurtailingParsers()
		var rs parser.ResultSet
		var err reader.Error
		for _, p := range parsers {
			h.RegisterCall()
			cp2, rs2, err2 := p.Parse(h, leftRecCtx, r.Clone())
			cp = cp.Union(cp2)
			rs.Append(rs2...)
			if err2 != nil && (err == nil || err2.Pos().Pos() >= err.Pos().Pos()) {
				err = err2
			}
			if len(rs2) > 0 {
				break
			}
		}
		if desc != "" && err != nil && err.Pos().Pos() == cur.Pos() {
			err = reader.NewError(cur, "was expecting %s", desc)
		}
		return cp, rs, err
	})
}
