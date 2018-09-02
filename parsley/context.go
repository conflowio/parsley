package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	reader      Reader
	resultCache ResultCache
	err         Error
	callCount   int
}

// NewContext creates a new parsing context
func NewContext(reader Reader) *Context {
	return &Context{
		reader:      reader,
		resultCache: NewResultCache(),
	}
}

// Reader returns with the reader
func (c *Context) Reader() Reader {
	return c.reader
}

// ResultCache returns with the result cache object
func (c *Context) ResultCache() ResultCache {
	return c.resultCache
}

// SetError saves the error if it has the highest position for found errors
func (c *Context) SetError(pos Pos, cause error) {
	if c.err == nil || int(c.err.Pos()) < int(pos) {
		c.err = NewError(pos, cause)
	}
}

// SetErrorf saves the error if it has the highest position for found errors
func (c *Context) SetErrorf(pos Pos, format string, values ...interface{}) {
	if c.err == nil || int(c.err.Pos()) < int(pos) {
		c.err = NewErrorf(pos, format, values...)
	}
}

// OverrideError overrides any previously set errors
func (c *Context) OverrideError(err Error) {
	c.err = err
}

// Error returns with the parse error with the highest position (if any)
func (c *Context) Error() Error {
	return c.err
}

// RegisterCall registers a call
func (c *Context) RegisterCall() {
	c.callCount++
}

// CallCount returns with the call count
func (c *Context) CallCount() int {
	return c.callCount
}
