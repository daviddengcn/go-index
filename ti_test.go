package index

import (
	"log"
	"github.com/daviddengcn/go-villa"
	"testing"
)

func TestTokenIndexer(t *testing.T) {
	outlinks := [][]string{
		{"a", /* -> */ "b", "c"},
		{"b", /* -> */ "a", "c"},
	}
	
	ti := &TokenIndexer{}
	
	for _, links := range outlinks {
		ti.Put(links[0], villa.NewStrSet(links[1:]...))
	}
	
	villa.AssertStringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[b]")
	villa.AssertStringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	villa.AssertStringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[a b]")
	
	var b villa.ByteSlice
	if err := ti.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
	}
	log.Printf("[ti] %d bytes written", len(b))

	if err := ti.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
	}
	
	villa.AssertStringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[b]")
	villa.AssertStringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	villa.AssertStringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[a b]")
	
	ti.Put("a", villa.NewStrSet("a", "b"))
	
	villa.AssertStringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[a b]")
	villa.AssertStringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	villa.AssertStringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[b]")
}
