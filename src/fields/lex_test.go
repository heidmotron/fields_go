package fields

import (
  "testing"
)

func collectItems(l *lexer) (is []item) {
  for i := range l.items {
    is = append(is, i)
  }
  return is
}

func TestIsAlphaUnder(t *testing.T) {
  r := isAlphaUnder('_')
  if !r {
    t.Errorf("underscore should be valid")
  }

  r = isAlphaUnder('Z')
  if !r {
    t.Errorf("Z should be valid")
  }
}

type lexTest struct {
  name     string
  input    string
  expected []item
}

var (
  tLeftParen  = item{itemLeftParen, 0, "("}
  tRightParen = item{itemRightParen, 0, ")"}
  tComma      = item{itemComma, 0, ","}

  tEOF = item{itemEOF, 0, ""}
)

var tests = []lexTest{
  {"fields and collections", "alpha,beta(yep)", []item{
    {itemField, 0, "alpha"},
    tComma,
    {itemField, 6, "beta"},
    tLeftParen,
    {itemField, 11, "yep"},
    tRightParen,
    tEOF,
  }},
  {"collection only", "b(c,d)", []item{
    {itemField, 0, "b"},
    tLeftParen,
    {itemField, 0, "c"},
    tComma,
    {itemField, 0, "d"},
    tRightParen,
    tEOF,
  }},
  {"invalid right paren", "a)", []item{
    {itemField, 0, "a"},
    {itemRightParen, 1, ")"},
    {itemError, 2, "unexpected ')'"},
  }},
}

func TestEmit(t *testing.T) {
  for _, test := range tests {
    l := lex(test.input)
    is := collectItems(l)

    for i, s := range is {
      if s.val != test.expected[i].val {
        t.Errorf("Test %s: Expected value %s got %s for %s", test.name, test.expected[i].val, s.val)
      }
    }
  }
}

func BenchmarkLex(b *testing.B) {
  for i := 0; i < b.N; i++ {
    l := lex("alpha,beta,common,yes(chicken(is),the(best(bargain))),whoa,goood")
    collectItems(l)
  }
}
