# Parsley - Parser combinator library written in Go

[![CircleCI status](https://circleci.com/gh/opsidian/parsley.svg?style=shield&circle-token=c42cce0e1ae1496645d1d6dc640d86a9e6de808d)](https://circleci.com/gh/opsidian/parsley) [![Codecov.io status](https://codecov.io/gh/opsidian/parsley/branch/master/graph/badge.svg)](https://codecov.io/gh/opsidian/parsley) [![GoDoc](https://godoc.org/github.com/opsidian/parsley?status.svg)](https://godoc.org/github.com/opsidian/parsley) [![Latest release](https://img.shields.io/github/release/opsidian/parsley.svg)](https://github.com/opsidian/parsley/releases/latest)

Parsley is a general parser combinator library which can be used to parse context-free, left-recursive languages. It handles indirect as well as direct left-recursion in polynomial time and defines a memoization helper for speeding up parsing time. The language grammar can be easily translated to a set of rules using parsers and combinators.

You can read a general introduction about parser combinators here: https://en.wikipedia.org/wiki/Parser_combinator

For more information about handling left-recursion please check out **Parser Combinators for Ambiguous Left-Recursive Grammars (2007)** by Frost R.A., Hafiz R., and Callaghan P.

Currently the library supports only text processing, but the interfaces are written with binary parsing in mind.

## How to use this library?

### Basics

#### Abstract Syntax Tree

An AST is built during parsing for evaluating the parsed data. An AST node has the following interface:

```
type Node interface {
	Token() string
	Value(ctx interface{}) (interface{}, Error)
	Pos() Pos
	ReaderPos() Pos
}
```

There are two types of nodes:

- **terminal node**: a leaf node - contains the smallest valid token, always has a constant value, e.g. a number, a string, etc.
- **non-terminal node**: a branch node - contains other nodes and an interpreter which defines how to evaluate its children together (example: the non-terminal node has the token "+", the terminal children are: [1, 2] and the interpreter would add the numbers and return 3)

The **ctx** evaluation context can be anything you would need for evaluating a tree. (e.g. looking up named variables in a variable store)

#### Interpreters

An interpreter gets the child nodes of a non-terminal node and returns with a single result. It has the following interface:

```
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, Error)
}
```

As you get the nodes as input you can implement lazy evaluation.

A simple example is when the interpreter function defines how to add numbers together.

#### Parsers

A parser has a simple interface:

```
type Parser interface {
	Parse(ctx Context, leftRecCtx data.IntMap, pos Pos) (Node, Error, data.IntSet)
	Name() string
}
```

A parser processes the next token(s) from the given reader and returns them as a result set. It also handles direct left recursion counters through leftRectCtx and accumulates the curtailing parsers in an int set. (You usually don't have to deal with these).

#### Combinators

Combinators are special parsers as they are combining other parsers to process more complex token groups. A simple example is the **Seq** combinator which simply tries to match the given parsers in order. Some combinator also use node builders which tells them how to build an AST node from the parsed token group.

#### Memoization and handling left-recursion

IMPORTANT: make sure you only call the Memoize generator function once for a specific parser as it generates an internal parser index for every call.

Depending on the language you define a parser can attempt to match at the same reader position multiple times. You can cache the parser results with the provided Memoize wrapper.

Also if your language contains left-recursion you need to use Memoize for any parser that is directly or indirectly part of it as Memoize is responsible for curtailing these calls.

The following code will wrap the integer parser with a memoizer:

```
p := combinator.Memoize(terminal.Integer())
```

#### A simple example

Let's write a parser which is able to parse the following expression: "INTEGER + INTEGER"

We'll need the following components:

- a parser which is able to match integer numbers ([terminal.Integer](text/terminal/integer.go))
- a parser which is able to match the "+" character ([terminal.Rune](text/terminal/rune.go))
- a combinator which is able to match multiple parsers in order ([combinator.Seq](combinator/seq.go))
- an interpreter function which will calculate the result

```
sum := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
	value0, _ := nodes[0].Value(ctx)
	value1, _ := nodes[2].Value(ctx)
	return value0.(int64) + value1.(int64), nil
})

p := combinator.SeqOf(
	terminal.Integer(),
	terminal.Rune('+'),
	terminal.Integer(),
).Bind(sum)

f := text.NewFile("example.file", []byte("1+2"))
fs := parsley.NewFileSet(f)
r := text.NewReader(f)
ctx := parsley.NewContext(fs, r)
value, err := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
if err != nil {
	panic(err)
}
fmt.Printf("Result: %d\n", value.(int64))
// Output: Result: 3
```

The **add** variable will contain a parser which is able to parse the given expression and the **s** variable contains the sentence parser which makes sure we parsed all the input data.

#### More examples

There is a JSON-like parser implementation in the [examples/json](examples/json) directory. It also gives you some

For a more complex expression parser you can check out the [Flint interpolation language](https://github.com/opsidian/flint).

### Documentation

Please more information about the available parsers and combinators please check out the [Go docs](https://godoc.org/github.com/opsidian/parsley).

## Library packages

- parsley (root): top level helper functions for parsing
- [ast](ast): abstract syntax tree related structs and interfaces
- [ast/interpreter](ast/interpreter): AST node interpreters
- [combinator](combinator): parser combinator implementations including memoization
- [data](data): int map and int set implementations
- [examples](examples): examples for how to use this library
- [parser](parser): the main parsing logic
- [parsley](parsley): common interfaces and the top-level parser/evaluate methods
- [text](text): text reader implementation
- [text/terminal](text/terminal): common parsers for text literals (string literal, int, float, etc.)

## Versioning

The library is expected to have occasional minor API changes while the version is 0.MAJOR.MINOR.

- MAJOR version will be increased for backward incompatible changes
- MINOR version will be increased for improvements and bugfixes

Starting from 1.0.0 the library will use the [Semantic Versioning 2.0.0](http://semver.org/spec/v2.0.0.html).

You can find the change log in the [CHANGELOG.md](CHANGELOG.md) file.

## Testing

To run all the tests simply run: `make test`.

## Benchmarking

There are some benchmarks for parsing JSON (which is obviously not the best comparison, but at least it's something).

```
$ cd examples/json/json
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/opsidian/parsley/examples/json/json
BenchmarkParsleyJSON1k-4      	   20000	     79929 ns/op	   26557 B/op	     570 allocs/op
BenchmarkParsleyJSON10k-4     	    2000	    681273 ns/op	  248000 B/op	    4435 allocs/op
BenchmarkParsleyJSON100k-4    	     200	   6999517 ns/op	 6011530 B/op	   41647 allocs/op
BenchmarkEncodingJSON1k-4     	   50000	     23834 ns/op	    5794 B/op	     113 allocs/op
BenchmarkEncodingJSON10k-4    	   10000	    190699 ns/op	   44140 B/op	     869 allocs/op
BenchmarkEncodingJSON100k-4   	    1000	   1880887 ns/op	  418492 B/op	    8153 allocs/op
PASS
ok  	github.com/opsidian/parsley/examples/json/json	11.410s
```

## LICENSE

This software is distributed under the Mozilla Public License, version 2.0. See the [LICENSE](LICENSE) file for more details.

## Acknowledgements

I would like to say thank you to Frost R.A., Hafiz R., and Callaghan P. as Parsley's main algorithms are based on their paper:

Frost R.A., Hafiz R., Callaghan P. (2007) Parser Combinators for Ambiguous Left-Recursive Grammars. In: Hudak P., Warren D.S. (eds) Practical Aspects of Declarative Languages. PADL 2008. Lecture Notes in Computer Science, vol 4902. Springer, Berlin, Heidelberg

---

Copyright (c) 2017 Opsidian Ltd.
