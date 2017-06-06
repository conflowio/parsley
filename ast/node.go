package ast

import (
	"fmt"

	"github.com/opsidian/parsley/reader"
)

// EOF is the end of file token
const EOF = "EOF"

// Node represents an AST node
type Node interface {
	Token() string
	Value() (interface{}, error)
}

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token string
	pos   reader.Position
	value interface{}
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, pos reader.Position, value interface{}) TerminalNode {
	return TerminalNode{
		token: token,
		pos:   pos,
		value: value,
	}
}

// Token returns with the node token
func (t TerminalNode) Token() string {
	return t.token
}

// Value returns with the value of the node
func (t TerminalNode) Value() (interface{}, error) {
	return t.value, nil
}

// Pos returns the position
func (t TerminalNode) Pos() reader.Position {
	return t.pos
}

func (t TerminalNode) String() string {
	if t.value == nil {
		return fmt.Sprintf("T{%s, %s}", t.token, t.pos)
	}
	return fmt.Sprintf("T{%v, %s}", t.value, t.pos)
}

// NonTerminalNode represents a non-leaf node in the AST
type NonTerminalNode struct {
	token       string
	children    []Node
	interpreter Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []Node, interpreter Interpreter) NonTerminalNode {
	return NonTerminalNode{
		token:       token,
		children:    children,
		interpreter: interpreter,
	}
}

// Token returns with the node token
func (n NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n NonTerminalNode) Value() (interface{}, error) {
	if n.interpreter == nil || len(n.children) == 0 {
		return nil, nil
	}
	nodes := make([]Node, 0, len(n.children))
	for _, child := range n.children {
		if child != nil {
			nodes = append(nodes, child)
		}
	}
	return n.interpreter.Eval(nodes)
}

// Children returns with the children
func (n NonTerminalNode) Children() []Node {
	return n.children
}

func (n NonTerminalNode) String() string {
	return fmt.Sprintf("NT{%s, %s}", n.token, n.children)
}
