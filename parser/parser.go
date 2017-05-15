package parser

import (
	"fmt"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// Parser defines a parser interface
type Parser interface {
	Parse(*History, *reader.Reader) Results
}

// Func defines a helper to implement the Parser interface with functions
type Func func(*History, *reader.Reader) Results

// Parse parses the next token and returns with an AST node and the updated reader
func (f Func) Parse(h *History, r *reader.Reader) Results {
	return f(h, r)
}

// Result represents one result of a parser
type Result struct {
	node   ast.Node
	reader *reader.Reader
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
	return fmt.Sprintf("%s, next: %d", r.node, r.reader.Position())
}

// Results is a result array
type Results []Result

// NewResults creates a new results instance
func NewResults(results ...Result) Results {
	return Results(results)
}

// Add adds a new result
func (r *Results) Add(node ast.Node, reader *reader.Reader) {
	*r = append(*r, Result{node, reader})
}

// History contains for each parser's last call position
type History struct {
	calls     map[string][]int
	callCount int
}

// NewHistory creates a new history instance
func NewHistory() *History {
	return &History{make(map[string][]int), 0}
}

// Push registers a parser call
func (h *History) Push(parser string, pos int) {
	if _, ok := h.calls[parser]; !ok {
		h.calls[parser] = []int{}
	}
	h.calls[parser] = append(h.calls[parser], pos)
	h.callCount++
}

// Pop removes a previous parser call
func (h *History) Pop(parser string) {
	h.calls[parser] = h.calls[parser][0 : len(h.calls[parser])-1]
}

// GetLastPosition returns with the last call position
func (h *History) GetLastPosition(parser string) (int, bool) {
	if calls, ok := h.calls[parser]; ok && len(h.calls[parser]) > 0 {
		return calls[len(calls)-1], true
	}
	return -1, false
}

// GetCallCount returns with the call count
func (h *History) GetCallCount() int {
	return h.callCount
}
