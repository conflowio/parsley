package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	fileSet     *FileSet
	reader      Reader
	resultCache ResultCache
	callCount   int
}

// NewContext creates a new parsing context
func NewContext(fileSet *FileSet, reader Reader) *Context {
	return &Context{
		fileSet:     fileSet,
		reader:      reader,
		resultCache: NewResultCache(),
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
