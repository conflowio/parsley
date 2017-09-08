package json

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text/terminal"
)

// NewParser returns with a new JSON parser
func NewParser() parser.Func {
	var value parser.Func

	array := combinator.Seq(
		builder.Select(1),
		terminal.Rune('[', "["),
		combinator.SepBy("ARRAY", &value, terminal.Rune(',', ","), arrayInterpreter()),
		terminal.Rune(']', "]"),
	)

	keyValue := combinator.Seq(builder.All("OBJ_KV", nil), terminal.String(), terminal.Rune(':', ":"), &value)

	object := combinator.Seq(
		builder.Select(1),
		terminal.Rune('{', "{"),
		combinator.SepBy("OBJ", keyValue, terminal.Rune(',', ","), objectInterpreter()),
		terminal.Rune('}', "}"),
	)

	value = combinator.Choice("value",
		terminal.String(),
		terminal.Float(),
		terminal.Integer(),
		array,
		object,
		terminal.Word("false", "FALSE", false),
		terminal.Word("true", "TRUE", true),
		terminal.Word("null", "NULL", nil),
	)

	return value
}

func arrayInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
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
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
		res := make(map[string]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			pair := nodes[i].(ast.NonTerminalNode)
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
