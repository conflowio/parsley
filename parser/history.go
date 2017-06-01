package parser

import (
	"github.com/opsidian/parsley/data"
)

type storedResult struct {
	curtailingParsers data.IntSet
	resultSet         ResultSet
	leftRecCtx        data.IntMap
}

// History records information about parser calls
type History struct {
	parserCount int
	parsers     map[string]int
	results     map[int]map[int]storedResult
}

// NewHistory creates a history instance
func NewHistory() *History {
	return &History{
		parserCount: 0,
		parsers:     make(map[string]int),
		results:     make(map[int]map[int]storedResult),
	}
}

// GetParserIndex maps the given parser to an integer index
func (h *History) GetParserIndex(parser string) (parserIndex int) {
	parserIndex, ok := h.parsers[parser]
	if !ok {
		parserIndex = h.parserCount
		h.parsers[parser] = parserIndex
		h.parserCount++
	}
	return
}

// RegisterResults registers a parser result for a certain position
func (h *History) RegisterResults(parserIndex int, pos int, curtailingParsers data.IntSet, resultSet ResultSet, leftRecCtx data.IntMap) {
	if _, ok := h.results[parserIndex]; !ok {
		h.results[parserIndex] = make(map[int]storedResult)
	}
	h.results[parserIndex][pos] = storedResult{curtailingParsers, resultSet, leftRecCtx}
}

// GetResults return with a previously saved result
func (h *History) GetResults(parserIndex int, pos int, leftRecCtx data.IntMap) (data.IntSet, ResultSet, bool) {
	storedResult, found := h.results[parserIndex][pos]
	if !found {
		return data.EmptyIntSet(), nil, false
	}

	for key := range storedResult.leftRecCtx.Keys() {
		if storedResult.leftRecCtx.Get(key) > leftRecCtx.Get(key) {
			return data.EmptyIntSet(), nil, false
		}
	}

	return storedResult.curtailingParsers, storedResult.resultSet, true
}
