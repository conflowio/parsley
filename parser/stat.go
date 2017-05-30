package parser

// Stat is the global statistics object
var Stat = Statistics{}

// Statistics contains parsing statistics
type Statistics struct {
	sumCallCount int
}

// RegisterCall registers a call
func (s *Statistics) RegisterCall() {
	s.sumCallCount++
}

// GetSumCallCount returns with the sum call count
func (s *Statistics) GetSumCallCount() int {
	return s.sumCallCount
}

// Reset resets the statistic counters
func (s *Statistics) Reset() {
	s.sumCallCount = 0
}
