// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Context is the parsing context passed to all parsers
type Context struct {
	fileSet               *FileSet
	reader                Reader
	resultCache           ResultCache
	err                   Error
	callCount             int
	keywords              map[string]struct{}
	transformationEnabled bool
	staticCheckEnabled    bool
	userCtx               interface{}
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

// EnableTransformation will turn on node transformation
func (c *Context) EnableTransformation() {
	c.transformationEnabled = true
}

// TransformationEnabled will return true if transformation is enabled
func (c *Context) TransformationEnabled() bool {
	return c.transformationEnabled
}

// EnableStaticCheck will turn on static checking
func (c *Context) EnableStaticCheck() {
	c.staticCheckEnabled = true
}

// StaticCheckEnabled will return true if static checking is enabled
func (c *Context) StaticCheckEnabled() bool {
	return c.staticCheckEnabled
}

// UserContext returns with the user context
func (c *Context) UserContext() interface{} {
	return c.userCtx
}

// SetUserContext sets the user context
func (c *Context) SetUserContext(userCtx interface{}) {
	c.userCtx = userCtx
}
