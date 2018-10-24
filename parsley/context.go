package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	fileSet     *FileSet
	reader      Reader
	resultCache ResultCache
	err         Error
	callCount   int
	keywords    map[string]struct{}
}

// NewContext creates a new parsing context
func NewContext(fileSet *FileSet, reader Reader) *Context {
	return &Context{
		fileSet:     fileSet,
		reader:      reader,
		resultCache: NewResultCache(),
		keywords:    make(map[string]struct{}, 64),
	}
}

// FileSet returns with the file set
func (c *Context) FileSet() *FileSet {
	return c.fileSet
}

// Reader returns with the reader
func (c *Context) Reader() Reader {
	return c.reader
}

// ResultCache returns with the result cache object
func (c *Context) ResultCache() ResultCache {
	return c.resultCache
}

// RegisterCall registers a call
func (c *Context) RegisterCall() {
	c.callCount++
}

// CallCount returns with the call count
func (c *Context) CallCount() int {
	return c.callCount
}

// SetError saves the error if it has the highest position for found errors
func (c *Context) SetError(err Error) {
	if err == nil {
		return
	}

	if c.err == nil || err.Pos() >= c.err.Pos() {
		c.err = err
	}
}

// Error returns with the parse error with the highest position (if any)
func (c *Context) Error() Error {
	return c.err
}

// RegisterKeywords registers one or more keywords
func (c *Context) RegisterKeywords(keywords ...string) {
	for _, keyword := range keywords {
		c.keywords[keyword] = struct{}{}
	}
}

// IsKeyword checks if the given string is a keyword
func (c *Context) IsKeyword(word string) bool {
	_, ok := c.keywords[word]
	return ok
}
