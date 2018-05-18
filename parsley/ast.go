package parsley

// Node represents an AST node
//go:generate counterfeiter . Node
type Node interface {
	Token() string
	Value(ctx interface{}) (interface{}, Error)
	Pos() Pos
	ReaderPos() Pos
}

// Interpreter defines an interface to evaluate AST nodes
//go:generate counterfeiter . Interpreter
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, Error)
}
