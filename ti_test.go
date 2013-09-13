package index

import (
	"testing"
	
	"github.com/daviddengcn/go-assert"
	"github.com/daviddengcn/go-villa"
)

func TestTokenIndexer(t *testing.T) {
	outlinks := [][]string{
		{"a" /* -> */, "b", "c"},
		{"b" /* -> */, "a", "c"},
	}

	ti := &TokenIndexer{}

	for _, links := range outlinks {
		ti.Put(links[0], villa.NewStrSet(links[1:]...))
	}

	assert.StringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[b]")
	assert.StringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	assert.StringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[a b]")

	var b villa.ByteSlice
	if err := ti.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
	}
	t.Logf("[ti] %d bytes written", len(b))

	if err := ti.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
	}

	assert.StringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[b]")
	assert.StringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	assert.StringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[a b]")

	ti.Put("a", villa.NewStrSet("a", "b"))

	assert.StringEquals(t, "inlinks of a", ti.IdsOfToken("a"), "[a b]")
	assert.StringEquals(t, "inlinks of b", ti.IdsOfToken("b"), "[a]")
	assert.StringEquals(t, "inlinks of c", ti.IdsOfToken("c"), "[b]")
}
