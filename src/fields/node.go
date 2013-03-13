package fields

import "fmt"

type Node interface {
  Type() NodeType
  String() string
}

type NodeType int

func (t NodeType) Type() NodeType {
  return t
}

const (
  NodeField NodeType = iota
  NodeCollection
  NodeList
)

type ListNode struct {
  NodeType
  Nodes []Node
}

func (l *ListNode) String() string {
  s := "\nLN\n"
  for _, node := range l.Nodes {
    s += "\t" + node.String() + "\n"
  }
  return s
}

type FieldNode struct {
  NodeType
  Text string
}

func (f *FieldNode) String() string {
  return f.Text
}

func newField(token item) *FieldNode {
  return &FieldNode{NodeField, token.val}
}

type CollectionNode struct {
  NodeType
  Field *FieldNode
  Nodes *ListNode
}

func (c *CollectionNode) String() string {
  return fmt.Sprintf("Collection Node: %s\n\t%s", c.Field.String(), c.Nodes.String())
}

func newCollection(field *FieldNode, nodes *ListNode) *CollectionNode {
  return &CollectionNode{NodeCollection, field, nodes}
}
