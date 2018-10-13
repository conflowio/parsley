package parsley

// Interpreter defines an interface to evaluate AST nodes
//go:generate counterfeiter . Interpreter
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, Error)
}
