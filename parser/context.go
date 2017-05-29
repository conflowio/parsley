package parser

import (
	"github.com/opsidian/parsec/data"
)

type storedResult struct {
	parserResult   *ParserResult
	leftRecContext data.IntMap
}

// Context records information about parser calls
type Context struct {
	parserCount  int
	sumCallCount int
	parsers      map[string]int
	results      map[int]map[int]storedResult
	depth        int
}

// NewContext creates a new context instance
func NewContext() *Context {
	return &Context{
		parserCount:  0,
		parsers:      make(map[string]int),
		results:      make(map[int]map[int]storedResult),
		sumCallCount: 1,
		depth:        0,
	}
}

// GetParserIndex maps the given parser to an integer index
func (c *Context) GetParserIndex(parser string) (parserIndex int) {
	parserIndex, ok := c.parsers[parser]
	if !ok {
		parserIndex = c.parserCount
		c.parsers[parser] = parserIndex
		c.parserCount++
	}
	return
}

// Log logs the given values with the current call stack depth
func (c *Context) Log(values ...interface{}) {
	//fmt.Print(strings.Repeat("  ", c.depth))
	//spew.Println(values...)
}

// RegisterCall registers a call
func (c *Context) RegisterCall() {
	c.sumCallCount++
	c.depth++
}

// FinishCall registers the end of a call
func (c *Context) FinishCall() {
	c.depth--
}

// GetSumCallCount returns with the sum call count
func (c *Context) GetSumCallCount() int {
	return c.sumCallCount
}

// RegisterResults registers a parser result for a certain position
func (c *Context) RegisterResults(parserIndex int, pos int, parserResult *ParserResult, leftRecContext data.IntMap) {
	if _, ok := c.results[parserIndex]; !ok {
		c.results[parserIndex] = make(map[int]storedResult)
	}
	c.results[parserIndex][pos] = storedResult{parserResult, leftRecContext}
}

// GetResults return with a previously saved result
func (c *Context) GetResults(parserIndex int, pos int, leftRecCtx data.IntMap) (*ParserResult, bool) {
	storedResult, found := c.results[parserIndex][pos]
	if !found {
		return nil, false
	}

	for key := range storedResult.leftRecContext.Keys() {
		if storedResult.leftRecContext.Get(key) > leftRecCtx.Get(key) {
			c.Log("RESNOUSE", storedResult, leftRecCtx)
			return nil, false
		}
	}

	c.Log("REUSE", storedResult.parserResult)
	return storedResult.parserResult, true
}
