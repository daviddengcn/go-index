package index

import (
	"testing"

	"github.com/golangplus/testing/assert"

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

	assert.Equal(t, "inlinks of a", ti.IdsOfToken("a"), []string{"b"})
	assert.Equal(t, "inlinks of b", ti.IdsOfToken("b"), []string{"a"})
	assert.Equal(t, "inlinks of c", ti.IdsOfToken("c"), []string{"a", "b"})

	// save/load
	var b villa.ByteSlice
	if err := ti.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
		return
	}
	t.Logf("[ti] %d bytes written", len(b))

	if err := ti.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
		return
	}

	assert.Equal(t, "inlinks of a", ti.IdsOfToken("a"), []string{"b"})
	assert.Equal(t, "inlinks of b", ti.IdsOfToken("b"), []string{"a"})
	assert.Equal(t, "inlinks of c", ti.IdsOfToken("c"), []string{"a", "b"})

	ti.Put("a", villa.NewStrSet("a", "b"))

	assert.Equal(t, "inlinks of a", ti.IdsOfToken("a"), []string{"a", "b"})
	assert.Equal(t, "inlinks of b", ti.IdsOfToken("b"), []string{"a"})
	assert.Equal(t, "inlinks of c", ti.IdsOfToken("c"), []string{"b"})
}
