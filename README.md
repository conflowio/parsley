# Parsley - Parser combinator library written in Go

[![CircleCI status](https://circleci.com/gh/opsidian/parsley.svg?style=shield&circle-token=c42cce0e1ae1496645d1d6dc640d86a9e6de808d)](https://circleci.com/gh/opsidian/parsley)

Parsley is a general parser combinator library which handles indirect as well as direct left-recursion in polynomial time and defines a memoization helper for speeding up parsing time.

You can read a general introduction about parser combinators here: https://en.wikipedia.org/wiki/Parser_combinator

For more information about handling left-recursion please check out **Parser Combinators for Ambiguous Left-Recursive Grammars (2007)** by Frost R.A., Hafiz R., and Callaghan P.

Currently only a text reader and text parsers are defined in the library, but the interfaces are written with  binary parsing in mind.

## How to use this library?

### Basics

#### Abstract Syntax Tree

An AST is built during parsing for evaluating the parsed data. An AST node has the following interface:

```
type Node interface {
	Token() string
	Value(ctx interface{}) (interface{}, error)
	Pos() reader.Position
}
```

There are two types of nodes:
 - **terminal node**: a leaf node - contains the smallest valid token, always has a constant value, e.g. a number, a string, etc.
 - **non-terminal node**: a branch node - contains other nodes and an interpreter which defines how to evaluate its children together (example: the non-terminal node has the token "+", the terminal children are: [1, 2] and the interpreter would add the numbers and return 3)

For an arithmetic expression like ```1 + 2 * 3``` we would build the following AST:

```
   +
 /   \
1     *
    /   \
   2     3
```

The **ctx** evaluation context can be anything you would need for evaluating a tree. (e.g. looking up named variables in a variable store)

#### Interpreters

An interpreter gets the child nodes of a non-terminal node and returns with a single result. It has the following interface:

```
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, error)
}
```

As you get the nodes as input you can implement lazy evaluation.

A simple example is when the interpreter function defines how to add numbers together.

#### Node builders

A node builder takes multiple nodes and returns with a new result node. It has the following simple interface:

```
type NodeBuilder interface {
	BuildNode([]Node) Node
}
```

An example is if the parser matches [1, +, 2] in order, the node builder would build the following non-terminal node with two terminal children:

```
   +
 /   \
1     2
```

#### Parsers

A parser has a simple interface:

```
type Parser interface {
	Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet, Error)
}
```

A parser processes the next token(s) from the given reader and returns them as a result set. It also handles direct left recursion counters through leftRectCtx and accumulates the curtailing parsers in an int set. (You usually don't have to deal with these).

One result is a tuple of an AST node and a reader clone with the new position.

#### Combinators

Combinators are special parsers as they are combining other parsers to process more complex token groups. A simple example is the **Seq** combinator which simply tries to match the given parsers in order. Some combinator also use node builders which tells them how to build an AST node from the parsed token group.

#### Memoization and handling left-recursion

Depending on the language you define a parser can attempt to match at the same reader position multiple times. You can cache the parser results with the provided Memoize wrapper.

Also if your language contains left-recursion you need to use Memoize for any parser that is directly or indirectly part of it as Memoize is responsible for curtailing these calls.

The following code will wrap the integer parser with a memoizer:

```
h := parser.NewHistory()
cachedParser := combinator.Memoize("INTEGER", h, terminal.Integer())
```

The history object will store the result cache and also track left recursion counts and curtailing parsers, so you should only create it once.

#### A simple example

Let's write a parser which is able to parse the following expression: "INTEGER + INTEGER"

We'll need the following components:
 - a parser which is able to match integer numbers ([terminal.Integer](text/terminal/integer.go))
 - a parser which is able to match the "+" character ([terminal.Rune](text/terminal/rune.go))
 - a combinator which is able to match multiple parsers in order ([combinator.Seq](combinator/seq.go))
 - a node builder function which will build the following AST tree: [+] with two children [number 1] and [number 2] ([builder.BinaryOperation](ast/builder/builder.go))
 - an interpreter function which will take two numeric nodes and will add them.
 - whitespaces will be ignored by the reader

```
add := combinator.Seq(
	builder.BinaryOperation(
		ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
			value0, _ := nodes[0].Value(ctx)
			value1, _ := nodes[1].Value(ctx)
			return value0.(int) + value1.(int), nil
		}),
	),
	terminal.Integer(),
	terminal.Rune('+', "ADD"),
	terminal.Integer(),
))
s := combinator.Seq(builder.Select(0), value, parser.End())

result, err := parsley.EvaluateText([]byte("1 + 2"), true, s, nil)
```

The **add** variable will contain a parser which is able to parse the given expression and the **s** variable contains the sentence parser which makes sure we parsed the full expression.

#### More examples

There is a JSON-like parser implementation in the [examples/json](examples/json) directory.

For a more complex expression parser you can check out the [Flint interpolation language](https://github.com/opsidian/flint).

### Documentation

Please more information about the available parsers and combinators please check out the [Go docs](https://godoc.org/github.com/opsidian/parsley).

## Library packages

 - parsley (root): top level helper functions for parsing
 - [ast](ast): abstract syntax tree related structs and interfaces
 - [ast/builder](ast/builder): AST node builder functions
 - [combinator](combinator): parser combinator implementations including memoization
 - [data](data): int map and int set implementations
 - [examples](examples): examples for how to use this library
 - [parser](parser): the main parsing logic
 - [reader](reader): interface definitions for an input reader
 - [text](text): text reader implementation
 - [text/terminal](text/terminal): common parsers for text literals (string literal, int, float, etc.)

## Acknowledgements

I would like to say thank you to Frost R.A., Hafiz R., and Callaghan P. as Parsley's main algorithms are based on their paper:

Frost R.A., Hafiz R., Callaghan P. (2007) Parser Combinators for Ambiguous Left-Recursive Grammars. In: Hudak P., Warren D.S. (eds) Practical Aspects of Declarative Languages. PADL 2008. Lecture Notes in Computer Science, vol 4902. Springer, Berlin, Heidelberg
