package fields

//import "fmt"
type Tree struct {
  Root      *ListNode
  lex       *lexer
  token     [2]item
  peekCount int
}

func (t *Tree) peek() item {
  if t.peekCount > 0 {
    return t.token[0]
  }

  t.peekCount = 1
  t.token[0] = t.lex.nextItem()
  return t.token[0]
}

func (t *Tree) next() item {
  if t.peekCount > 0 {
    t.peekCount--
  } else {
    t.token[0] = t.lex.nextItem()
  }

  return t.token[t.peekCount]
}

func (t *Tree) backup() {
  t.peekCount++
}

func (t *Tree) parse() {
  for t.peek().typ != itemEOF {
    node := t.parseFieldOrCollection()
    t.Root.Nodes = append(t.Root.Nodes, node)
  }
}

/**
  (collection / field) *
*/
func (t *Tree) parseFieldOrCollection() Node {
  token := t.next()
  if token.typ == itemField {
    nextToken := t.peek()
    switch nextToken.typ {
    case itemComma:
      t.next()
    case itemLeftParen:
      t.next()
      collection := newCollection(newField(token), t.parseCollection())
      if t.peek().typ == itemComma {
        t.next()
      }
      return collection
    }
    return newField(token)
  }
  panic("where am i")
  return nil
}

func (t *Tree) parseCollection() *ListNode {
  list := &ListNode{NodeType: NodeList}
  list.Nodes = t.parseInner()
  var token item
  if token = t.next(); token.typ != itemRightParen {
    panic("No closing right paren")
  }
  return list
}

func (t *Tree) parseInner() []Node {
  var nodes []Node

  for t.peek().typ != itemRightParen {
    nodes = append(nodes, t.parseFieldOrCollection())
  }
  return nodes
}

func Parse(input string) (*Tree, error) {
  t := &Tree{
    Root: &ListNode{NodeType: NodeList},
    lex:  lex(input),
  }

  t.parse()
  return t, nil
}
