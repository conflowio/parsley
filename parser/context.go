package parser

type StoredResult struct {
	parserResult   *ParserResult
	leftRecContext IntMap
}

// Context records information about parser calls
type Context struct {
	parserCount  int
	sumCallCount int
	parsers      map[string]int
	results      map[int]map[int]StoredResult
}

// NewContext creates a new context instance
func NewContext() *Context {
	return &Context{
		parserCount:  0,
		parsers:      make(map[string]int),
		results:      make(map[int]map[int]StoredResult),
		sumCallCount: 1,
	}
}

func (c *Context) GetParserIndex(parser string) (parserIndex int) {
	parserIndex, ok := c.parsers[parser]
	if !ok {
		parserIndex = c.parserCount
		c.parsers[parser] = parserIndex
		c.parserCount++
	}
	return
}

// RegisterCall registers a call
func (c *Context) RegisterCall() {
	c.sumCallCount++
}

// GetSumCallCount returns with the sum call count
func (c *Context) GetSumCallCount() int {
	return c.sumCallCount
}

// RegisterResults registers a parser result for a certain position
func (c *Context) RegisterResults(parserIndex int, pos int, parserResult *ParserResult, leftRecContext IntMap) {
	if _, ok := c.results[parserIndex]; !ok {
		c.results[parserIndex] = make(map[int]StoredResult)
	}
	c.results[parserIndex][pos] = StoredResult{parserResult, leftRecContext}
}

// GetResults return with a previously saved result
func (c *Context) GetResults(parserIndex int, pos int, leftRecCtx IntMap) (*ParserResult, bool) {
	storedResult, found := c.results[parserIndex][pos]
	if !found {
		return nil, false
	}

	if storedResult.leftRecContext != nil {
		for srParserIndex, srLeftRecCount := range storedResult.leftRecContext {
			if srLeftRecCount > leftRecCtx[srParserIndex] {
				return nil, false
			}
		}
	}

	return storedResult.parserResult, true
}
