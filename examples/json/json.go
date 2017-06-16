package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opsidian/parsley"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text/terminal"
)

func main() {
	jsonFilePath := "example.json"
	if len(os.Args) > 1 {
		jsonFilePath = os.Args[1]
	}
	b, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}

	h := parser.NewHistory()

	var value parser.Func

	array := combinator.Seq(
		builder.Select(1),
		terminal.Rune('[', "["),
		combinator.SepBy("ARRAY", h, &value, terminal.Rune(',', ","), arrayInterpreter()),
		terminal.Rune(']', "]"),
	)

	keyValue := combinator.Seq(builder.All("OBJ_KV", nil), terminal.String(), terminal.Rune(':', ":"), &value)

	object := combinator.Seq(
		builder.Select(1),
		terminal.Rune('{', "{"),
		combinator.SepBy("OBJ", h, keyValue, terminal.Rune(',', ","), objectInterpreter()),
		terminal.Rune('}', "}"),
	)

	value = combinator.Memoize("VALUE", h, combinator.Any(
		terminal.String(),
		terminal.Integer(),
		terminal.Float(),
		array,
		object,
		terminal.Word("false", "FALSE", false),
		terminal.Word("true", "TRUE", true),
		terminal.Word("null", "NULL", nil),
	))

	s := combinator.Seq(builder.Select(0), object, parser.End())
	res, err := parsley.EvaluateText(b, true, s)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parser calls: %d\n", parser.Stat.GetSumCallCount())
	fmt.Println(string(json))
}

func arrayInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
		res := make([]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			value, err := nodes[i].Value()
			if err != nil {
				return nil, err
			}
			res[i/2] = value
		}
		return res, nil
	})
}

func objectInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
		res := make(map[string]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			pair := nodes[i].(ast.NonTerminalNode)
			key, _ := pair.Children()[0].Value()
			value, err := pair.Children()[2].Value()
			if err != nil {
				return nil, err
			}
			res[key.(string)] = value
		}
		return res, nil
	})
}
