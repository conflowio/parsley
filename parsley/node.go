package parsley

// Node represents an AST node
//go:generate counterfeiter . Node
type Node interface {
	Token() string
	Type() string
	Value(ctx interface{}) (interface{}, Error)
	Pos() Pos
	ReaderPos() Pos
}

// NonTerminalNode represents a nonterminal AST node
//go:generate counterfeiter . NonTerminalNode
type NonTerminalNode interface {
	Node
	Children() []Node
	StaticCheck(ctx interface{}) Error
}
