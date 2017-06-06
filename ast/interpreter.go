package ast

// Interpreter defines an interface to evaluate an AST node
type Interpreter interface {
	Eval([]Node) (interface{}, error)
}

// InterpreterFunc defines a helper to implement the Interpreter interface with functions
type InterpreterFunc func([]Node) (interface{}, error)

// Eval evaluates the given values and returns with a single result
func (f InterpreterFunc) Eval(nodes []Node) (interface{}, error) {
	return f(nodes)
}
