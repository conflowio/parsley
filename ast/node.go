package ast

import (
	"fmt"

	"github.com/opsidian/parsec/reader"
)

// Node represents an AST node
type Node interface {
	Token() string
	Value() (interface{}, error)
	Pos() reader.Position
}

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token string
	pos   reader.Position
	value interface{}
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, pos reader.Position, value interface{}) *TerminalNode {
	return &TerminalNode{
		token: token,
		pos:   pos,
		value: value,
	}
}

// Token returns with the node token
func (t *TerminalNode) Token() string {
	return t.token
}

// Value returns with the value of the node
func (t *TerminalNode) Value() (interface{}, error) {
	return t.value, nil
}

// Pos returns the position
func (t *TerminalNode) Pos() reader.Position {
	return t.pos
}

func (t *TerminalNode) String() string {
	if t.value != nil {
		return fmt.Sprintf("T{%v, %s}", t.value, t.pos)
	} else {
		return fmt.Sprintf("T{%s, %s}", t.token, t.pos)
	}
}

// Interpreter is a function to evaluate an AST node
type Interpreter func(children []interface{}) (interface{}, error)

// NonTerminalNode represents a non-leaf node in the AST
type NonTerminalNode struct {
	token       string
	children    []Node
	interpreter Interpreter
	pos         reader.Position
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []Node, interpreter Interpreter) *NonTerminalNode {
	return &NonTerminalNode{
		token:       token,
		children:    children,
		interpreter: interpreter,
		pos:         children[0].Pos(),
	}
}

// Token returns with the node token
func (n *NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value() (interface{}, error) {
	values := make([]interface{}, len(n.children))
	for i, child := range n.children {
		value, err := child.Value()
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return n.interpreter(values)
}

// Pos returns the position
func (n *NonTerminalNode) Pos() reader.Position {
	return n.pos
}

func (n *NonTerminalNode) String() string {
	return fmt.Sprintf("N{%s, %s}", n.token, n.children)
}
