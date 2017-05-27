package parser

import (
	"fmt"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/data"
	"github.com/opsidian/parsec/reader"
)

// Result represents one result of a parser
type Result struct {
	node   ast.Node
	reader *reader.Reader
}

// NewResult creates a new result instance
func NewResult(node ast.Node, reader *reader.Reader) Result {
	return Result{
		node:   node,
		reader: reader,
	}
}

// Node returns with the node
func (r Result) Node() ast.Node {
	return r.node
}

// Reader returns with the reader
func (r Result) Reader() *reader.Reader {
	return r.reader
}

func (r Result) String() string {
	return fmt.Sprintf("%s, cur: %s", r.node, r.reader.Cursor())
}

// ParserResult is the result of a parse call
type ParserResult struct {
	CurtailingParsers data.IntSet
	Results           []Result
}

// NewParserResult creates a new parser result
func NewParserResult(curtailingParsers data.IntSet, results ...Result) *ParserResult {
	return &ParserResult{curtailingParsers, results}
}

// Append adds a result to the parse result
func (p *ParserResult) Append(results ...Result) {
	for _, result := range results {
		p.append(result)
	}
}

func (p *ParserResult) append(result Result) {
	if p.Results == nil {
		p.Results = []Result{result}
		return
	}

	for k, v := range p.Results {
		if v.Reader().Cursor().Pos() > result.Reader().Cursor().Pos() {
			p.Results = append(p.Results, Result{})
			copy(p.Results[k+1:], p.Results[k:])
			p.Results[k] = result
			return
		}
	}

	p.Results = append(p.Results, result)
}
