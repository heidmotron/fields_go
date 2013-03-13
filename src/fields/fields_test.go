package fields

import "testing"

func TestParse(t *testing.T) {
  // _, err := Parse("")
  // if err == nil {
  //   t.Errorf("Failed to raise an error on empty string")
  // }
  //
  tree, err := Parse("a")
  if err != nil {
    t.Errorf("Should not raise an error")
  }

  if len(tree.Root.Nodes) != 1 {
    t.Errorf("Should only have one node")
  }

  if tree.Root.Nodes[0].String() != "a" {
    t.Errorf("Should be Node Field")
  }

  tree, err = Parse("a,b")

  if tree != nil && tree.Root.Nodes[1].String() != "b" {
     t.Errorf("failed to parse single field")
  }

//  tree, err = Parse("b(c,d),d(e,f(g))")
//  t.Errorf("%v", tree)

//  if tree != nil && tree.Root.Nodes[0].String() != "b" {
//    t.Errorf("failed to parse single field, %s", tree.Root.Nodes[0].String())
//  }
}

func BenchmarkParse(b *testing.B) {
  for i:=0; i <  b.N; i++ {
    Parse("alpha,beta,common,yes(chicken(is),the(best(bargain))),whoa,goood")
  }
}
