// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// This is a JSON parser example. It's not a strict implementation as it was written only for demonstration. For most JSON strings it should still be able to parse the input.
//
// You can run this file to see the parser in action:
//  go run json.go
// By default the included example.json file will be used and the output will be:
//  Parser calls: 240
//  map[title:Person type:object properties:map[firstName:map[type:string] lastName:map[type:string] age:map[description:Age in years type:integer minimum:0]] required:[firstName lastName]]
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
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

	value = combinator.Memoize(combinator.Any("value",
		terminal.String(),
		terminal.Integer(),
		terminal.Float(),
		array,
		object,
		terminal.Word("false", "FALSE", false),
		terminal.Word("true", "TRUE", true),
		terminal.Word("null", "NULL", nil),
	))

	s := parsley.NewSentence(object)
	res, h, err := s.Evaluate(text.NewReader(b, true), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parser calls: %d\n", h.CallCount())
	fmt.Printf("%v\n", res)
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
