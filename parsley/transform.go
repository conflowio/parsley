package parsley

// Transformable defines an interface about a transformable Node
type Transformable interface {
	Transform(userCtx interface{}) (Node, Error)
}

// TransformableNode defines a transformable node
//counterfeiter:generate . TransformableNode
type TransformableNode interface {
	Node
	Transformable
}

// NodeTransformer defines an interface to transform an AST node to an other
//counterfeiter:generate . NodeTransformer
type NodeTransformer interface {
	TransformNode(userCtx interface{}, node Node) (Node, Error)
}

// NodeTransformFunc is a function which implements the NodeTransformer interface
type NodeTransformFunc func(userCtx interface{}, node Node) (Node, Error)

// TransformNode transforms an AST node to an other
func (f NodeTransformFunc) TransformNode(userCtx interface{}, node Node) (Node, Error) {
	return f(userCtx, node)
}

// NodeTransformerRegistry contains named node registries
type NodeTransformerRegistry interface {
	NodeTransformer(name string) (NodeTransformer, bool)
}

// Transform transforms the given node recursively, in a breadth-first manner
func Transform(userCtx interface{}, node Node) (Node, Error) {
	switch n := node.(type) {
	case Transformable:
		return n.Transform(userCtx)
	}

	return node, nil
}
