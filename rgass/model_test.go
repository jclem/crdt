package rgass_test

import (
	"testing"

	"github.com/jclem/crdt/rgass"
)

func TestFindNode(t *testing.T) {
	m := rgass.NewModel()
	id := rgass.ID{Length: 10}
	node := &rgass.Node{
		ID:    id,
		Str:   "1234567890",
		Split: true,
		List: []*rgass.Node{
			&rgass.Node{ID: rgass.ID{Length: 2}, Str: "12"},
			&rgass.Node{ID: rgass.ID{Offset: 2, Length: 5}, Split: true, Str: "34567", List: []*rgass.Node{
				&rgass.Node{Str: "3", ID: rgass.ID{Length: 1}},
				&rgass.Node{Str: "4567", ID: rgass.ID{Offset: 1, Length: 4}},
			}},
			&rgass.Node{ID: rgass.ID{Offset: 7, Length: 3}, Split: true, Str: "890", List: []*rgass.Node{
				&rgass.Node{Str: "89", ID: rgass.ID{Length: 2}},
				&rgass.Node{Str: "0", ID: rgass.ID{Offset: 2, Length: 1}},
			}},
		},
	}
	m.InsertAfter(m.Head(), node)
	if n, _ := m.FindNode(id, 0); n.Str != "12" {
		t.Fatalf("Expected %s, got %s", "12", n.Str)
	}
	if n, _ := m.FindNode(id, 2); n.Str != "12" {
		t.Fatalf("Expected %s, got %s", "12", n.Str)
	}
	if n, _ := m.FindNode(id, 3); n.Str != "3" {
		t.Fatalf("Expected %s, got %s", "3", n.Str)
	}
	if n, _ := m.FindNode(id, 10); n.Str != "0" {
		t.Fatalf("Expected %s, got %s", "0", n.Str)
	}
	if _, err := m.FindNode(id, 11); err == nil {
		t.Fatalf("Expected an error, got none")
	}
}
