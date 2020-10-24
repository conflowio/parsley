package parsley

// Walkable is a generic interface to allow to apply a function on the node
// The Walk function should return true if the walk should be interrupted
type Walkable interface {
	Walk(f func(n Node) bool) bool
}

// WalkableNode defines a node which also implements the Walkable interface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . WalkableNode
type WalkableNode interface {
	Node
	Walk(f func(n Node) bool) bool
}

// Walk applies the given function to the node recursively, in a depth-first manner
func Walk(node Node, f func(n Node) bool) bool {
	switch n := node.(type) {
	case Walkable:
		if n.Walk(f) {
			return true
		}
	case NonTerminalNode:
		for _, child := range n.Children() {
			if Walk(child, f) {
				return true
			}
		}
	}

	return f(node)
}
