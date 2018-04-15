package json

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() *parser.NamedFunc {
	var value parser.NamedFunc

	array := combinator.Seq("ARRAY",
		terminal.Rune('['),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Array()),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	keyValue := combinator.Seq("OBJ_KV",
		terminal.String(false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(&value, text.WsSpaces),
	)

	object := combinator.Seq("OBJ",
		terminal.Rune('{'),
		combinator.SepBy(
			text.LeftTrim(keyValue, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Object()),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	value = *combinator.Choice("value",
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Word("false", false),
		terminal.Word("true", true),
		terminal.Word("null", nil),
	)

	return text.RightTrim(&value, text.WsSpacesNl)
}
