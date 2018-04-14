package json

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() parser.Func {
	var value parser.NamedFunc

	array := combinator.Seq("ARRAY",
		terminal.Rune('['),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Bind(arrayInterpreter())

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
		),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Bind(objectInterpreter())

	value = *combinator.Memoize(combinator.Choice("value",
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Word("false", false),
		terminal.Word("true", true),
		terminal.Word("null", nil),
	))

	return text.RightTrim(&value, text.WsSpacesNl)
}

func arrayInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		if nodes[1].Token() == ast.NIL {
			return []interface{}{}, nil
		}

		nodes = nodes[1].(*ast.NonTerminalNode).Children()
		res := make([]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			value, err := nodes[i].Value(ctx)
			if err != nil {
				return nil, err
			}
			res[i/2] = value
		}
		return res, nil
	})
}

func objectInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		if nodes[1].Token() == ast.NIL {
			return map[string]interface{}{}, nil
		}

		nodes = nodes[1].(*ast.NonTerminalNode).Children()
		res := make(map[string]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			pair := nodes[i].(*ast.NonTerminalNode)
			key, _ := pair.Children()[0].Value(ctx)
			value, err := pair.Children()[2].Value(ctx)
			if err != nil {
				return nil, err
			}
			res[key.(string)] = value
		}
		return res, nil
	})
}
