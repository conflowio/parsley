package parsley

// Interpreter defines an interface to evaluate the given nonterminal node
//go:generate counterfeiter . Interpreter
type Interpreter interface {
	Eval(ctx interface{}, node NonTerminalNode) (interface{}, Error)
}

// StaticCheckerInterpreter defines an interpreter which is also a static checker
//go:generate counterfeiter . StaticCheckerInterpreter
type StaticCheckerInterpreter interface {
	Interpreter
	StaticChecker
}
