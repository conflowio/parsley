package parser

import (
	"fmt"

	"github.com/opsidian/parsec/ast"
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

// AsList wraps the result as an array
func (r Result) AsList() *Results {
	return NewResults(nil, r)
}

// Results contains a result array and reasons for curtailing if empty
type Results struct {
	items          []Result
	curtailReasons map[int]int
}

// NewResults creates a new results instance
func NewResults(curtailReasons map[int]int, items ...Result) *Results {
	if items == nil {
		items = make([]Result, 0)
	}
	r := &Results{
		items:          items,
		curtailReasons: nil,
	}
	if curtailReasons != nil {
		r.curtailReasons = make(map[int]int, len(curtailReasons))
		for k, v := range curtailReasons {
			r.curtailReasons[k] = v
		}
	}
	return r
}

// Items returns with the result items
func (r *Results) Items() []Result {
	return r.items
}

// Add adds a new result
func (r *Results) Add(result Result) {
	r.items = append(r.items, result)
}

// Merge merges two result list
func (r *Results) Merge(r2 *Results) {
	r.items = append(r.items, r2.items...)
	r.MergeCurtailReasons(r2)
}

// MergeCurtailReasons merges the curtail reason maps
func (r *Results) MergeCurtailReasons(results *Results) {
	if results == nil || results.curtailReasons == nil {
		return
	}
	if r.curtailReasons == nil {
		r.curtailReasons = make(map[int]int)
	}
	for k, v := range results.curtailReasons {
		if v > r.curtailReasons[k] {
			r.curtailReasons[k] = v
		}
	}
}

func (r *Results) String() string {
	return fmt.Sprintf("Results{%v}", r.items)
}
