package parser_test

/*func TestRegisterCallShouldIncreaseSumCallCount(t *testing.T) {
	c := parser.NewContext()
	c.RegisterCall()
	assert.Equal(t, c.GetSumCallCount(), 2)
	c.RegisterCall()
	assert.Equal(t, c.GetSumCallCount(), 3)
}

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	c := parser.NewContext()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(nil, parser.NewResult(node, nil))
	c.RegisterResults(c.GetParserIndex("p1"), 2, results, parser.NewIntMap())

	actual, ok := c.GetResults(c.GetParserIndex("p1"), 2, parser.NewIntMap())
	assert.Equal(t, results, actual)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	c := parser.NewContext()
	c.RegisterResults(c.GetParserIndex("p1"), 2, nil, parser.NewIntMap())
	results, ok := c.GetResults(c.GetParserIndex("p1"), 2, parser.NewIntMap())
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	c := parser.NewContext()
	results, ok := c.GetResults(c.GetParserIndex("p1"), 2, parser.NewIntMap())
	assert.Nil(t, results)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	c := parser.NewContext()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(nil, parser.NewResult(node, nil))
	c.RegisterResults(c.GetParserIndex("p1"), 1, results, parser.NewIntMap())
	c.RegisterResults(c.GetParserIndex("p2"), 2, nil, parser.NewIntMap())

	actual, ok := c.GetResults(c.GetParserIndex("p1"), 1, parser.NewIntMap())
	assert.Equal(t, results, actual)
	assert.True(t, ok)

	results, ok = c.GetResults(c.GetParserIndex("p2"), 2, parser.NewIntMap())
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestGetSumCallCountShouldStartAtOne(t *testing.T) {
	c := parser.NewContext()
	assert.Equal(t, c.GetSumCallCount(), 1)
}

func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
	c := parser.NewContext()
	ctx := parser.NewIntMap()
	ctx[c.GetParserIndex("p1")] = 2
	ctx[c.GetParserIndex("p2")] = 1
	c.RegisterResults(c.GetParserIndex("p1"), 1, parser.NewParserResult(ctx), parser.NewIntMap())

	ctx = parser.NewIntMap()
	ctx[c.GetParserIndex("p1")] = 1
	ctx[c.GetParserIndex("p2")] = 1
	ctx[c.GetParserIndex("p3")] = 1
	results, found := c.GetResults(c.GetParserIndex("p1"), 1, parser.NewIntMap())
	assert.Nil(t, results)
	assert.False(t, found)
}

func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
	c := parser.NewContext()
	ctx := parser.NewIntMap()
	ctx[c.GetParserIndex("p1")] = 2
	ctx[c.GetParserIndex("p2")] = 1
	cResults := parser.NewParserResult(ctx, parser.NewResult(nil, nil))
	c.RegisterResults(c.GetParserIndex("p1"), 1, cResults, parser.NewIntMap())
	c.RegisterResults(c.GetParserIndex("p2"), 1, cResults, parser.NewIntMap())

	ctx = parser.NewIntMap()
	ctx[c.GetParserIndex("p1")] = 1
	ctx[c.GetParserIndex("p2")] = 1

	var found bool
	ctx[c.GetParserIndex("p3")] = 1
	_, found = c.GetResults(c.GetParserIndex("p3"), 1, parser.NewIntMap())
	assert.False(t, found)

	ctx[c.GetParserIndex("p1")] = 1
	_, found = c.GetResults(c.GetParserIndex("p1"), 1, parser.NewIntMap())
	assert.False(t, found)

	ctx[c.GetParserIndex("p2")] = 1
	_, found = c.GetResults(c.GetParserIndex("p2"), 1, parser.NewIntMap())
	assert.False(t, found)

	ctx[c.GetParserIndex("p1")] = 1
	results, found := c.GetResults(c.GetParserIndex("p1"), 1, parser.NewIntMap())
	assert.Equal(t, cResults, results)
	assert.True(t, found)
}*/
