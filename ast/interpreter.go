package ast

// Interpreter defines an interface to evaluate AST nodes
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, error)
}

// InterpreterFunc defines a helper to implement the Interpreter interface with functions
type InterpreterFunc func(ctx interface{}, nodes []Node) (interface{}, error)

// Eval evaluates the given nodes and returns with a single result.
func (f InterpreterFunc) Eval(ctx interface{}, nodes []Node) (interface{}, error) {
	return f(ctx, nodes)
}
