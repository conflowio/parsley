// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json

import (
	"github.com/conflowio/parsley/ast/interpreter"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() parsley.Parser {
	var value parser.Func

	array := combinator.SeqOf(
		terminal.Rune('['),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Array()),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	keyValue := combinator.SeqOf(
		terminal.String("string", false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(&value, text.WsSpacesNl),
	)

	object := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.SepBy(
			text.LeftTrim(keyValue, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Object()),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	value = combinator.Choice(
		terminal.String("string", false),
		terminal.Float("number"),
		terminal.Integer("integer"),
		array,
		object,
		terminal.Bool("boolean", "true", "false"),
		terminal.Nil("null", "null"),
	).Name("value")

	return value
}
