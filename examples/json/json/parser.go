// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() parsley.Parser {
	var value parser.Func

	array := combinator.SeqOf(
		terminal.Rune('['),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Transform(combinator.TransformNodes("ARRAY", interpreter.Array(), false)),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Transform(combinator.TransformNodes("OBJECT", interpreter.Select(1), false))

	keyValue := combinator.SeqOf(
		terminal.String(false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(&value, text.WsSpacesNl),
	)

	object := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.SepBy(
			text.LeftTrim(keyValue, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Transform(combinator.TransformNodes("OBJECT", interpreter.Object(), false)),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Transform(combinator.TransformNodes("OBJECT", interpreter.Select(1), false))

	value = combinator.Choice(
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Bool("true", "false"),
		terminal.Nil("null"),
	).Name("value")

	return value
}
