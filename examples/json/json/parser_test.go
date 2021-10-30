// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json_test

import (
	"io/ioutil"
	"testing"

	encoding_json "encoding/json"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/examples/json/json"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

func benchmarkParsleyJSON(b *testing.B, jsonFilePath string) {
	f, err := text.ReadFile(jsonFilePath)
	if err != nil {
		b.Fatal(err)
	}

	s := combinator.Sentence(text.Trim(json.NewParser()))
	r := text.NewReader(f)
	ctx := parsley.NewContext(parsley.NewFileSet(f), r)
	ctx.EnableStaticCheck()
	ctx.EnableTransformation()
	if _, err = parsley.Evaluate(ctx, s); err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		_, _ = parsley.Evaluate(ctx, s)
	}
}

func BenchmarkParsleyJSON1k(b *testing.B)   { benchmarkParsleyJSON(b, "../example_1k.json") }
func BenchmarkParsleyJSON10k(b *testing.B)  { benchmarkParsleyJSON(b, "../example_10k.json") }
func BenchmarkParsleyJSON100k(b *testing.B) { benchmarkParsleyJSON(b, "../example_100k.json") }

func benchmarkEncodingJSON(b *testing.B, jsonFilePath string) {
	input, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		b.Fatal(err)
	}

	var val interface{}
	if err := encoding_json.Unmarshal(input, &val); err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		var val interface{}
		_ = encoding_json.Unmarshal(input, &val)
	}
}

func BenchmarkEncodingJSON1k(b *testing.B)   { benchmarkEncodingJSON(b, "../example_1k.json") }
func BenchmarkEncodingJSON10k(b *testing.B)  { benchmarkEncodingJSON(b, "../example_10k.json") }
func BenchmarkEncodingJSON100k(b *testing.B) { benchmarkEncodingJSON(b, "../example_100k.json") }
