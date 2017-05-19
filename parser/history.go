package parser

// History contains each parser's last call position
type History struct {
	parserCounter int
	parsers       map[string]int
	calls         map[int]map[int]int
	results       map[int]map[int]Results
	callCount     int
}

// NewHistory creates a new history instance
func NewHistory() *History {
	return &History{
		parserCounter: 1,
		parsers:       make(map[string]int),
		calls:         make(map[int]map[int]int),
		callCount:     0,
		results:       make(map[int]map[int]Results),
	}
}

// Push registers a parser call
func (h *History) Push(parser string, pos int) {
	parserIndex, ok := h.parsers[parser]
	if !ok {
		h.parsers[parser] = h.parserCounter
		parserIndex = h.parserCounter
		h.parserCounter++
	}
	if _, ok := h.calls[parserIndex]; !ok {
		h.calls[parserIndex] = make(map[int]int)
	}
	if _, ok := h.calls[parserIndex][pos]; !ok {
		h.calls[parserIndex][pos] = 1
	} else {
		h.calls[parserIndex][pos]++
	}
	h.callCount++
}

// Pop removes a previous parser call
func (h *History) Pop(parser string, pos int) {
	h.calls[h.parsers[parser]][pos]--
}

// GetCalls returns with the call count for a given position
func (h *History) GetCalls(parser string, pos int) int {
	return h.calls[h.parsers[parser]][pos]
}

// GetCallCount returns with the call count
func (h *History) GetCallCount() int {
	return h.callCount
}

// RegisterResults registers a parser result for a certain position
func (h *History) RegisterResults(parser string, pos int, results Results) {
	parserIndex, ok := h.parsers[parser]
	if !ok {
		h.parsers[parser] = h.parserCounter
		parserIndex = h.parserCounter
		h.parserCounter++
	}
	if _, ok := h.results[parserIndex]; !ok {
		h.results[parserIndex] = make(map[int]Results)
	}
	h.results[parserIndex][pos] = results
}

// GetResults return with a previously saved result
func (h *History) GetResults(parser string, pos int) (results Results, ok bool) {
	parserIndex, ok := h.parsers[parser]
	if !ok {
		return nil, false
	}
	results, ok = h.results[parserIndex][pos]
	return
}
