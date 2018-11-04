package parsley

// StaticChecker defines an interface to run a static analysis on the given nonterminal node
//go:generate counterfeiter . StaticChecker
type StaticChecker interface {
	StaticCheck(ctx interface{}, node NonTerminalNode) (string, Error)
}
