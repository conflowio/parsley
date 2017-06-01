package parser

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/reader"
)

// Result represents one result of a parser
type Result struct {
	node   ast.Node
	reader reader.Reader
}

// NewResult creates a new result instance
func NewResult(node ast.Node, reader reader.Reader) Result {
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
func (r Result) Reader() reader.Reader {
	return r.reader
}

func (r Result) String() string {
	return fmt.Sprintf("RES{%s, cur: %s}", r.node, r.reader.Cursor())
}

// AsSet converts the result to a result set containing this result
func (r Result) AsSet() ResultSet {
	return ResultSet([]Result{r})
}

// ResultSet is a set of results
type ResultSet []Result

// NewResultSet creates a new result set
func NewResultSet(results ...Result) ResultSet {
	return ResultSet(results)
}

// Append adds a result to the parse result
func (rs *ResultSet) Append(results ...Result) {
	for _, result := range results {
		rs.append(result)
	}
}

func (rs *ResultSet) append(result Result) {
	if *rs == nil {
		*rs = []Result{result}
		return
	}

	for k, v := range *rs {
		// If we already have a result up to the same position then we ignore it
		if v.Reader().Cursor().Pos() == result.Reader().Cursor().Pos() {
			return
		}
		if v.Reader().Cursor().Pos() > result.Reader().Cursor().Pos() {
			*rs = append(*rs, Result{})
			copy((*rs)[k+1:], (*rs)[k:])
			(*rs)[k] = result
			return
		}
	}

	*rs = append(*rs, result)
}
