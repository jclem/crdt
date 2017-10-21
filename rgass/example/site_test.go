package example_test

import (
	"testing"

	"github.com/jclem/crdt/rgass/example"
)

func TestSiteInsert(t *testing.T) {
	site := example.NewSite(1, 1)

	if err := site.Insert(0, "Hello, world"); err != nil {
		t.Fatalf(err.Error())
	}
	if text := site.Text(); text != "Hello, world" {
		t.Fatalf("Expected %q, got %q", "Hello, world", text)
	}
	if err := site.Insert(5, " there"); err != nil {
		t.Fatalf(err.Error())
	}
	if text := site.Text(); text != "Hello there, world" {
		t.Fatalf("Expected %q, got %q", "Hello there, world", text)
	}
	if err := site.Insert(8, "!"); err != nil {
		t.Fatalf(err.Error())
	}
	if text := site.Text(); text != "Hello th!ere, world" {
		t.Fatalf("Expected %q, got %q", "Hello th!ere, world", text)
	}

	op := <-site.OutStream
	if str := op.Str; str != "Hello, world" {
		t.Fatalf("Expected %q, got %q", "Hello, world", str)
	}
}

func TestSiteDelete(t *testing.T) {
	site := example.NewSite(1, 1)

	if err := site.Insert(0, "Hello, world"); err != nil {
		t.Fatalf(err.Error())
	}
	if err := site.Insert(5, " th..ere"); err != nil {
		t.Fatalf(err.Error())
	}
	if err := site.Delete(8, 2); err != nil {
		t.Fatalf(err.Error())
	}
	if text := site.Text(); text != "Hello there, world" {
		t.Fatalf("Expected %q, got %q", "Hello there, world", text)
	}

	op := <-site.OutStream
	op = <-site.OutStream
	op = <-site.OutStream
	if len := op.Len; len != 2 {
		t.Fatalf("Expected %q, got %q", 2, len)
	}
}

func TestConcurrent(t *testing.T) {
	site1 := example.NewSite(1, 1)
	site2 := example.NewSite(1, 2)

	// Site 1 inserts "Hello, world"
	if err := site1.Insert(0, "Helloworld"); err != nil {
		t.Fatalf(err.Error())
	}

	// Site 2 incorporates s1op1
	if err := site2.Receive(<-site1.OutStream); err != nil {
		t.Fatalf(err.Error())
	}

	// Site 2 inserts into "Hello, world"
	if err := site2.Insert(5, ", "); err != nil {
		t.Fatalf(err.Error())
	}

	// Site 1 concurrently deletes "lowo"
	if err := site1.Delete(3, 4); err != nil {
		t.Fatalf(err.Error())
	}

	// Site 2 incorporates s1op2
	if err := site2.Receive(<-site1.OutStream); err != nil {
		t.Fatalf(err.Error())
	}

	if text := site2.Text(); text != "Hel, rld" {
		t.Fatalf("Expected %q, got %q", "Hel, rld", text)
	}

	// Site 1 incorporates s2op1
	if err := site1.Receive(<-site2.OutStream); err != nil {
		t.Fatalf(err.Error())
	}

	if site1.Text() != site2.Text() {
		t.Fatalf("Site 1 had %q, site 2 had %q", site1.Text(), site2.Text())
	}
}
