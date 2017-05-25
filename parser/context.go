package parser

// Context records information about parser calls
type Context struct {
	parserCount    int
	sumCallCount   int
	parsers        map[string]int
	callCounts     map[int]map[int]int
	results        map[int]map[int]*Results
	curtailReasons map[int]map[int]map[int]int
}

// NewContext creates a new context instance
func NewContext() *Context {
	return &Context{
		parserCount:    0,
		parsers:        make(map[string]int),
		callCounts:     make(map[int]map[int]int),
		results:        make(map[int]map[int]*Results),
		curtailReasons: make(map[int]map[int]map[int]int),
		sumCallCount:   1,
	}
}

func (c *Context) getParserIndex(parser string) (parserIndex int) {
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

// Push registers a parser call
func (c *Context) Push(parser string, pos int) {
	parserIndex := c.getParserIndex(parser)
	if _, ok := c.callCounts[pos]; !ok {
		c.callCounts[pos] = make(map[int]int)
	}
	if _, ok := c.callCounts[pos][parserIndex]; !ok {
		c.callCounts[pos][parserIndex] = 1
	} else {
		c.callCounts[pos][parserIndex]++
	}
}

// Pop removes a previous parser call
func (c *Context) Pop(parser string, pos int) {
	parserIndex := c.getParserIndex(parser)
	if c.callCounts[pos][parserIndex] == 0 {
		panic("Pop was called without a matching Push")
	}
	c.callCounts[pos][parserIndex]--
}

// GetCalls returns with the call count for a given position
func (c *Context) GetCalls(parser string, pos int) int {
	parserIndex := c.getParserIndex(parser)
	return c.callCounts[pos][parserIndex]
}

// GetSumCallCount returns with the sum call count
func (c *Context) GetSumCallCount() int {
	return c.sumCallCount
}

// RegisterResults registers a parser result for a certain position
func (c *Context) RegisterResults(parser string, pos int, results *Results) {
	if _, ok := c.results[pos]; !ok {
		c.results[pos] = make(map[int]*Results)
	}
	c.results[pos][c.getParserIndex(parser)] = results
}

// GetResults return with a previously saved result
func (c *Context) GetResults(parser string, pos int) (results *Results, found bool) {
	parserIndex := c.getParserIndex(parser)
	results, found = c.results[pos][parserIndex]
	if found && results != nil {
		for crParserIndex, crCallCount := range results.curtailReasons {
			if crCallCount > c.callCounts[pos][crParserIndex] {
				return nil, false
			}
		}
	}
	return
}

// NewCurtailedResults creates a new empty result with the curtailing reasons
func (c *Context) NewCurtailedResults(pos int) *Results {
	return NewResults(c.callCounts[pos])
}
