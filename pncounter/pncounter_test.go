package pncounter_test

import (
	"testing"

	"github.com/jclem/crdt/pncounter"
)

func TestIncrement(t *testing.T) {
	g := pncounter.NewPNCounter("A")
	g.Increment()
	if v := g.Value(); v != 1 {
		t.Fatalf("Expected 1, got %q", v)
	}
}

func TestIncorporate(t *testing.T) {
	g := pncounter.NewPNCounter("A")
	g.Increment()
	g.Increment()
	g.Decrement()
	g.Incorporate("B", [2]int{2, 3})
	g.Incorporate("B", [2]int{1, 4})
	g.Incorporate("C", [2]int{3, 0})

	if exp, v := 2, g.Value(); v != 2 {
		t.Fatalf("Expected %d, got %d", exp, v)
	}
}
