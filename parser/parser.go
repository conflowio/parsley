package parser

import (
	"errors"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Parser defines a parser interface
type Parser interface {
	Parse(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult
}

// Func defines a helper to implement the Parser interface with functions
type Func func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult

// Parse parses the next token and returns with an AST node and the updated reader
func (f Func) Parse(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
	return f(leftRecCtx, r)
}

// Parse parses the given input with the parser function
func Parse(input []byte, p Func) (interface{}, error) {
	r := reader.New(input, true)
	parserResult := p.Parse(data.NewIntMap(), r)
	if len(parserResult.Results) == 0 {
		return nil, errors.New("Failed to parse the input")
	}
	return parserResult.Results[0].Node().Value()
}
