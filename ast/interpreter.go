package ast

// Interpreter defines an interface to evaluate an AST node
type Interpreter interface {
	Eval([]interface{}) (interface{}, error)
}

// InterpreterFunc defines a helper to implement the Interpreter interface with functions
type InterpreterFunc func([]interface{}) (interface{}, error)

// Eval evaluates the given values and returns with a single result
func (f InterpreterFunc) Eval(values []interface{}) (interface{}, error) {
	return f(values)
}
