package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	reader  Reader
	history History
	err     Error
}

// NewContext creates a new parsing context
func NewContext(reader Reader, history History) *Context {
	return &Context{
		reader:  reader,
		history: history,
	}
}

// Reader returns with the reader
func (c *Context) Reader() Reader {
	return c.reader
}

// History returns with the history object
func (c *Context) History() History {
	return c.history
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
