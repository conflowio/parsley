// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// String matches a string literal enclosed in double quotes
func String() parser.Func {
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		tr := r.(*text.Reader)
		matches, pos, ok := tr.ReadMatch("\"|`", false)
		if !ok {
			return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting string literal")
		}
		quote := matches[0]

		var value string
		if quote == "`" {
			var matches []string
			matches, _, _ = tr.ReadMatch("[^`]*", true)
			value = matches[0]
		} else {
			value, _, _ = tr.Readf(unquoteString, true)
		}

		endQuote, _, err := tr.ReadRune()
		if err != nil || string(endQuote) != quote {
			return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting '%s'", quote)
		}
		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("STRING", pos, value), tr).AsSet(), nil
	})
}

func unquoteString(b []byte) (string, int, bool) {
	str := string(b)
	var tail, res string
	var err error
	var ch rune
	for {
		if str == "" {
			break
		}
		ch, _, tail, err = strconv.UnquoteChar(str, '"')
		if err != nil {
			break
		}
		res += string(ch)
		str = tail
	}
	return res, len(b) - len(str), true
}
