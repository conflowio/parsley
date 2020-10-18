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
	"os"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/examples/json/json"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func main() {
	jsonFilePath := "example.json"
	if len(os.Args) > 1 {
		jsonFilePath = os.Args[1]
	}
	file, err := text.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}
	fs := parsley.NewFileSet(file)

	reader := text.NewReader(file)
	ctx := parsley.NewContext(fs, reader)
	s := combinator.Sentence(text.Trim(json.NewParser()))

	res, evalErr := parsley.Evaluate(ctx, s)
	if evalErr != nil {
		panic(evalErr)
	}
	fmt.Printf("Parser calls: %d\n", ctx.CallCount())
	fmt.Printf("%v\n", res)
}
