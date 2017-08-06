# Parsley - Parser combinator library written in Go

You can read a general introduction about parser combinators here: https://en.wikipedia.org/wiki/Parser_combinator

## How to use this library?

The best place to start is to check out a lazy JSON parser implementation in the examples directory.

For a more complex expression parser you can check out the [Flint interpolation language](https://github.com/opsidian/flint).

## Library packages

 - parsley (root): top level helper functions for parsing
 - [ast](ast): abstract syntax tree related structs and interfaces
 - [ast/builder](ast/builder): abstract syntax tree builder functions
 - [combinator](combinator): parser combinator implementations including memoization
 - [data](data): int map and int set implementations
 - [examples](examples): examples for how to use this library
 - [parser](parser): the main parsing logic
 - [reader](reader): interface definitions for an input reader
 - [text](text): text reader implementation
 - [text/terminal](text/terminal): common parsers for text literals (string literal, int, float, etc.)

## Acknowledgements

The parser combinator is based on the following paper:

Frost R.A., Hafiz R., Callaghan P. (2007) Parser Combinators for Ambiguous Left-Recursive Grammars. In: Hudak P., Warren D.S. (eds) Practical Aspects of Declarative Languages. PADL 2008. Lecture Notes in Computer Science, vol 4902. Springer, Berlin, Heidelberg
