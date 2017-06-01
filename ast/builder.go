package ast

// NodeBuilder defines an interface to build a node from multiple nodes
type NodeBuilder interface {
	BuildNode([]Node) Node
}

// NodeBuilderFunc defines a helper to implement the NodeBuilder interface with functions
type NodeBuilderFunc func([]Node) Node

// BuildNode combines the given nodes into a single node
func (f NodeBuilderFunc) BuildNode(nodes []Node) Node {
	return f(nodes)
}
