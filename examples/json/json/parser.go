package json

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() parsley.ParserFunc {
	var value parsley.ParserFunc

	array := combinator.Seq(
		builder.All("ARRAY", interpreter.Select(1)),
		terminal.Rune('['),
		combinator.SepBy("ARRAY_ELEMENTS", &value, terminal.Rune(','), arrayInterpreter()),
		terminal.Rune(']'),
	)

	keyValue := combinator.Seq(builder.All("OBJ_KV", nil), terminal.String(false), terminal.Rune(':'), &value)

	object := combinator.Seq(
		builder.All("OBJ", interpreter.Select(1)),
		terminal.Rune('{'),
		combinator.SepBy("OBJ_ATTRIBUTES", keyValue, terminal.Rune(','), objectInterpreter()),
		terminal.Rune('}'),
	)

	value = combinator.Choice("value",
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Word("FALSE", "false", false),
		terminal.Word("TRUE", "true", true),
		terminal.Word("NULL", "null", nil),
	)

	return value
}

func arrayInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
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
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
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
