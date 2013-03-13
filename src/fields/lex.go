package fields

import (
  "fmt"
  "strings"
  "unicode/utf8"
)

const eof = -1

type itemType int

const (
  itemError itemType = iota
  itemField
  itemLeftParen
  itemRightParen
  itemComma
  itemEOF
)

type item struct {
  typ itemType
  pos int
  val string
}

type stateFn func (*lexer) stateFn

type lexer struct {
  input string
  start int
  pos int
  width int
  lastPos int
  parenDepth int
  state stateFn
  items chan item
}

func lex(input string) *lexer {
  l := &lexer{
    input: input,
    items: make(chan item),
  }

  go l.run()
  return l
}

func (l *lexer) run() {
  for l.state = lexField; l.state != nil; {
     l.state = l.state(l)
   }
   close(l.items)
}

func (l *lexer) emit(t itemType) {
  l.items <- item{typ: t, pos: l.start, val: l.input[l.start:l.pos]}
  l.start = l.pos
}

func (l *lexer) nextItem() item  {
  return <- l.items
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
  msg := fmt.Sprintf(format, args...)
  l.items <- item{typ: itemError, pos: l.pos, val: msg}
  return nil
}

func (l *lexer) next() rune {
  if l.pos >= len(l.input) {
    l.width = 0
    return eof
  }

  r, w := utf8.DecodeRuneInString(l.input[l.pos:])
  l.width = w
  l.pos += l.width

  return r
}

func (l *lexer) ignore() {
  l.start = l.pos
}

func (l *lexer) backup() {
  l.pos -= l.width
}

func (l *lexer) peek() rune {
  r := l.next()
  l.backup()
  return r
}

// State functions

func lexField(l *lexer) stateFn {
  for {
    rest := l.input[l.pos:]
    if strings.HasPrefix(rest, ",") {
      if (l.pos > l.start) {
        l.emit(itemField)
      }
      return lexComma
    }

    if strings.HasPrefix(rest, "(") {
      if (l.pos > l.start) {
        l.emit(itemField)
        l.parenDepth++
      }
      return lexLeftParen
    }

    if strings.HasPrefix(rest, ")") {
      if (l.pos > l.start) {
        l.emit(itemField)
      }
      return lexRightParen
    }

    r := l.next()
    if r == eof  {
      break
    } else if !isAlphaUnder(r) {
      return l.errorf("unexpected '%s'", r)
    }
  }

  if l.pos > l.start {
    l.emit(itemField)
  }
  l.emit(itemEOF)
  return nil
}

func lexComma(l *lexer) stateFn {
  l.pos += len(",")
  l.emit(itemComma)

  return lexField
}

func lexLeftParen(l *lexer) stateFn {
  l.pos += len("(")
  l.emit(itemLeftParen)
  return lexField
}

func lexRightParen(l *lexer) stateFn {
  l.pos += len(")")
  l.emit(itemRightParen)
  l.parenDepth--
  if (l.parenDepth < 0) {
    return l.errorf("unexpected ')'")
  }

  return lexField
}

func isAlphaUnder(r rune) bool {
  return r == '_' || (r >= 'A' && r <= 'z')
}

