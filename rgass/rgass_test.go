package rgass_test

import (
	"fmt"
	"testing"

	"github.com/jclem/crdt/rgass"
)

func TestLocalDeleteWhole(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if _, _, err := rg.LocalDelete(id, 0, 4); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "" {
		t.Fatalf("Expected %q, got: %q", "", text)
	}
}

func TestLocalDeletePrior(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if _, _, err := rg.LocalDelete(id, 0, 2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "st" {
		t.Fatalf("Expected %q, got: %q", "st", text)
	}
}

func TestLocalDeleteLast(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if _, _, err := rg.LocalDelete(id, 2, 2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "te" {
		t.Fatalf("Expected %q, got: %q", "te", text)
	}
}

func TestLocalDeleteMiddle(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if _, _, err := rg.LocalDelete(id, 1, 2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "tt" {
		t.Fatalf("Expected %q, got: %q", "tt", text)
	}
}

func TestLocalDeleteMultiple(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "1234", id); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	id2 := site.NextID(5)
	if err := rg.LocalInsert(id, 4, "!@#$%", id2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	id3 := site.NextID(5)
	if err := rg.LocalInsert(id2, 5, "abcde", id3); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if _, _, err := rg.LocalDelete(id, 2, 9); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "12cde" {
		t.Fatalf("Expected %q, got: %q", "12cde", text)
	}
}

func TestLocalInsertHead(t *testing.T) {
	// Test an insert at the head
	site := Site{}
	rg := rgass.NewRGASS()
	err := rg.LocalInsert(rg.Head().ID, 0, "test", site.NextID(4))
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "test" {
		t.Fatalf("Expected %q, got: %q", "test", text)
	}
}

func TestLocalInsertEnd(t *testing.T) {
	// Test an insert at the end of a node
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.LocalInsert(id1, 4, "test2", site.NextID(5)); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "testtest2" {
		t.Fatalf("Expected %q, got: %q", "testtest2", text)
	}
}

func TestLocalInsertInside(t *testing.T) {
	// Test an insert inside of a node
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	err := rg.LocalInsert(rg.Head().ID, 0, "test", id1)
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	err = rg.LocalInsert(id1, 2, "test2", site.NextID(5))
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if text := rg.Text(); text != "tetest2st" {
		t.Fatalf("Expected %q, got: %q", "tetest2st", text)
	}
}

func TestLocalInsertMultiple(t *testing.T) {
	// Test an insert that needs to come after a prior insert
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	id2 := site.NextID(5)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test2", id2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if text := rg.Text(); text != "test2test" {
		t.Fatalf("Expected %q, got: %q", "test2test", text)
	}
}

func TestRemoteInsert(t *testing.T) {
	TestLocalInsertHead(t)
	TestLocalInsertEnd(t)
	TestLocalInsertInside(t)
	TestLocalInsertMultiple(t)
}

func TestRemoteDeleteSingleNonSplit(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.RemoteDelete([]rgass.ID{id1}, 1, 2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if text := rg.Text(); text != "tt" {
		t.Fatalf("Expected %q, got: %q", "tt", text)
	}
}

func TestRemoteDeleteSingleSplit(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "test", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id2 := site.NextID(4)
	if err := rg.LocalInsert(id1, 2, "1234", id2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.RemoteDelete([]rgass.ID{id1}, 1, 2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if text := rg.Text(); text != "t1234t" {
		t.Fatalf("Expected %q, got: %q", "t1234t", text)
	}
}

func TestRemoteDeleteMultipleNonSplit(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "1234", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id2 := site.NextID(4)
	if err := rg.LocalInsert(id1, 4, "5678", id2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id3 := site.NextID(4)
	if err := rg.LocalInsert(id2, 4, "90ab", id3); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.RemoteDelete([]rgass.ID{id1, id2, id3}, 2, 8); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if text := rg.Text(); text != "12ab" {
		t.Fatalf("Expected %q, got: %q", "12ab", text)
	}
}

func TestRemoteDeleteMultipleSplit(t *testing.T) {
	site := Site{}
	rg := rgass.NewRGASS()
	id1 := site.NextID(4)
	if err := rg.LocalInsert(rg.Head().ID, 0, "1234", id1); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id2 := site.NextID(4)
	if err := rg.LocalInsert(id1, 2, "xxxx", id2); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id1Child := rg.MustGet(id1).List[1].ID
	id3 := site.NextID(4)
	if err := rg.LocalInsert(id1Child, 2, "5678", id3); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id4 := site.NextID(4)
	if err := rg.LocalInsert(id3, 2, "yyyy", id4); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id3Child := rg.MustGet(id3).List[1].ID
	id5 := site.NextID(4)
	if err := rg.LocalInsert(id3Child, 2, "90ab", id5); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	id6 := site.NextID(4)
	if err := rg.LocalInsert(id5, 2, "zzzz", id6); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if err := rg.RemoteDelete([]rgass.ID{id1Child, rg.MustGet(id3).List[0].ID, rg.MustGet(id3).List[1].ID, rg.MustGet(id5).List[0].ID}, 0, 8); err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
	if text := rg.Text(); text != "12xxxxyyyyzzzzab" {
		t.Fatalf("Expected %q, got: %q", "12xxxxyyyyzzzzab", text)
	}
}

func ExampleRGASS_Head() {
	rg := rgass.NewRGASS()
	id := rgass.ID{Vector: 1, Length: 5}
	if err := rg.LocalInsert(rg.Head().ID, 0, "Hello", id); err != nil {
		panic(err)
	}
	fmt.Println(rg.Text())
	// Output: Hello
}

func ExampleRGASS_LocalInsert() {
	rg := rgass.NewRGASS()
	id := rgass.ID{Vector: 1, Length: 5}
	if err := rg.LocalInsert(rg.Head().ID, 0, "Hello", id); err != nil {
		panic(err)
	}
	id2 := rgass.ID{Vector: 2, Length: 5}
	if err := rg.LocalInsert(id, 2, "WORLD", id2); err != nil {
		panic(err)
	}
	fmt.Println(rg.Text())
	// Output: HeWORLDllo
}

func ExampleRGASS_LocalDelete() {
	rg := rgass.NewRGASS()
	id := rgass.ID{Vector: 1, Length: 5}
	if err := rg.LocalInsert(rg.Head().ID, 0, "Hello", id); err != nil {
		panic(err)
	}
	if _, _, err := rg.LocalDelete(id, 1, 3); err != nil {
		panic(err)
	}
	fmt.Println(rg.Text())
	// Output: Ho
}

type Site struct {
	Vector int
}

func (s *Site) NextID(length int) rgass.ID {
	s.Vector++
	return rgass.ID{Length: length, Offset: 0, Session: 0, Site: 0, Vector: s.Vector}
}
