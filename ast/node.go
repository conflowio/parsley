package ast

import "go/token"

// Interpreter is a function to evaluate an
type Interpreter func(children []interface{}) (interface{}, error)

// Node represents an AST node
type Node interface {
	Token() token.Token
	Value() (interface{}, error)
}

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token token.Token
	value interface{}
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token token.Token, value interface{}) *TerminalNode {
	return &TerminalNode{
		token: token,
		value: value,
	}
}

// Token returns with the node token
func (t *TerminalNode) Token() token.Token {
	return t.token
}

// Value returns with the value of the node
func (t *TerminalNode) Value() (interface{}, error) {
	return t.value, nil
}

// NonTerminalNode represents a non-leaf node in the AST
type NonTerminalNode struct {
	token       token.Token
	children    []Node
	interpreter Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token token.Token, children []Node, interpreter Interpreter) *NonTerminalNode {
	return &NonTerminalNode{
		token:       token,
		children:    children,
		interpreter: interpreter,
	}
}

// Token returns with the node token
func (n *NonTerminalNode) Token() token.Token {
	return n.token
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value() (interface{}, error) {
	childrenValues := make([]interface{}, len(n.children))
	for i, child := range n.children {
		childValue, err := child.Value()
		if err != nil {
			return nil, err
		}
		childrenValues[i] = childValue
	}
	return n.interpreter(childrenValues)
}
