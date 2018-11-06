// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	fileSet             *FileSet
	reader              Reader
	resultCache         ResultCache
	err                 Error
	callCount           int
	keywords            map[string]struct{}
	nodeFactoryRegistry NodeFactoryRegistry
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

// SetNodeFactoryRegistry sets the node factory registry
func (c *Context) SetNodeFactoryRegistry(nodeFactoryRegistry NodeFactoryRegistry) {
	c.nodeFactoryRegistry = nodeFactoryRegistry
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

// NodeFactoryExists returns true if the named factory exists
func (c *Context) NodeFactoryExists(factoryName string) bool {
	if c.nodeFactoryRegistry == nil {
		return false
	}
	return c.nodeFactoryRegistry.NodeFactoryExists(factoryName)
}

// GetNodeFactory returns with the named factory if it exists
func (c *Context) GetNodeFactory(factoryName string) (NodeFactory, bool) {
	if c.nodeFactoryRegistry == nil {
		return nil, false
	}

	return c.nodeFactoryRegistry.GetNodeFactory(factoryName)
}
