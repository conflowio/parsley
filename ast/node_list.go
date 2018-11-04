package ast

import "github.com/opsidian/parsley/parsley"

// NodeList contains a list of nodes, should be used when a parser returns with multiple results
// If you call any of the parsley.Node methods on it then it will behave as it would be the first node.
type NodeList []parsley.Node

// Token returns with NODE_LIST
func (nl NodeList) Token() string {
	return nl[0].Token()
}

// Type returns with the type of the first node
func (nl NodeList) Type() string {
	return nl[0].Type()
}

// Value returns with the value of the first node
func (nl NodeList) Value(ctx interface{}) (interface{}, parsley.Error) {
	return nl[0].Value(ctx)
}

// Pos returns the value of the first node
func (nl NodeList) Pos() parsley.Pos {
	return nl[0].Pos()
}

// ReaderPos returns the reader position for the first node
func (nl NodeList) ReaderPos() parsley.Pos {
	return nl[0].ReaderPos()
}

// SetReaderPos amends the reader position for all nodes
func (nl NodeList) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	for i, node := range nl {
		nl[i] = SetReaderPos(node, f)
	}
}

// Walk runs the given function on the first node
func (nl NodeList) Walk(f func(n parsley.Node) bool) bool {
	return WalkNode(nl[0], f)
}

// Append appends a new node to the list
func (nl *NodeList) Append(node parsley.Node) {
	switch v := node.(type) {
	case NodeList:
		for _, node := range v {
			nl.Append(node)
		}
	case EmptyNode:
		for _, node := range *nl {
			if node == v {
				return
			}
		}
		*nl = append(*nl, v)
	default:
		*nl = append(*nl, v)
	}
}
