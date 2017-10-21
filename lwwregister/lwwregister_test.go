package lwwregister_test

import "testing"
import "github.com/jclem/crdt/lwwregister"

func TestUpdate(t *testing.T) {
	r := lwwregister.NewRegister(1)
	r.Update(1)
	if v := r.Val.(int); v != 1 {
		t.Fatalf("Expected 1, got %d", v)
	}
}

func TestIncorporate(t *testing.T) {
	r := lwwregister.NewRegister(1)
	r.Update(1)
	ts1 := lwwregister.Timestamp{2, 2}
	r.Incorporate(ts1, 2)
	if v := r.Val.(int); v != 2 {
		t.Fatalf("Expected 2, got %d", v)
	}
	ts2 := lwwregister.Timestamp{3, 0}
	r.Incorporate(ts2, 3)
	if v := r.Val.(int); v != 2 {
		t.Fatalf("Expected 2, got %d", v)
	}
}
