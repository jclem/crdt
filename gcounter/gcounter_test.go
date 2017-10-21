package gcounter_test

import (
	"testing"

	"github.com/jclem/crdt/gcounter"
)

func TestIncrement(t *testing.T) {
	g := gcounter.NewGCounter("A")
	g.Increment()
	if v := g.Value(); v != 1 {
		t.Fatalf("Expected 1, got %q", v)
	}
}

func TestIncorporate(t *testing.T) {
	g := gcounter.NewGCounter("A")
	g.Increment()
	g.Incorporate("B", 2)
	g.Incorporate("B", 1)
	g.Incorporate("C", 3)

	if v := g.Value(); v != 6 {
		t.Fatalf("Expected 6, got %q", v)
	}
}
