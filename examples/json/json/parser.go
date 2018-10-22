package json

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() *combinator.Sequence {
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
		terminal.String(false),
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
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Word("false", false),
		terminal.Word("true", true),
		terminal.Word("null", nil),
	).Name("value")

	return combinator.SeqOf(
		text.LeftTrim(value, text.WsSpacesNl),
		terminal.Whitespaces(text.WsSpacesNl),
	).Bind(interpreter.Select(0))
}
