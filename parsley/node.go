package parsley

// Node represents an AST node
//go:generate counterfeiter . Node
type Node interface {
	Token() string
	Value(ctx interface{}) (interface{}, Error)
	Pos() Pos
	ReaderPos() Pos
}
