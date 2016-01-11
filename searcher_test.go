package index

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/golangplus/bytes"
	"github.com/golangplus/strings"
	"github.com/golangplus/testing/assert"
)

type DocInfo struct {
	A string
}

func indexDocs(docs [][2]string) *TokenSetSearcher {
	sch := &TokenSetSearcher{}
	for i := range docs {
		text := docs[i][1]
		var tokens stringsp.Set
		TokenizeBySeparators(" ,", bytesp.NewPSlice([]byte(text)),
			func(token []byte) error {
				tokens.Add(string(token))
				return nil
			})

		fields := map[string]stringsp.Set{
			"text": tokens,
		}
		sch.AddDoc(fields, &DocInfo{
			A: fmt.Sprintf("%d - %s", i+1, docs[i][0]),
		})
	}

	return sch
}

func TestTokenSetSearcher(t *testing.T) {
	DOCS := [][2]string{
		{"To friends", "hello my friend"},
		{"To dogs", "GO go go, my dog"},
	}

	sch := indexDocs(DOCS)

	var docs []int32
	var infos []*DocInfo
	collector := func(docID int32, data interface{}) error {
		docs = append(docs, docID)
		docInfo := data.(*DocInfo)
		infos = append(infos, docInfo)
		t.Logf("Doc: %d, %+v\n", docID, docInfo)
		return nil
	}
	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "my"), collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	assert.Equal(t, "len(docs)(my)", len(docs), 2)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "my", "dog"), collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	assert.Equal(t, "len(docs)(my dog)", len(docs), 1)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "friend"), collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	assert.Equal(t, "len(docs)(friend)", len(docs), 1)

	docs = sch.TokenDocList("text", "my")
	assert.StringEqual(t, "text:To", docs, "[0 1]")

	sch.Delete(0)

	docs, infos = nil, nil
	sch.Search(nil, collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	assert.Equal(t, "len(docs)()", len(docs), 1)

	var b bytesp.Slice
	gob.Register(&DocInfo{})
	if err := sch.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
		return
	}
	t.Logf("%d bytes written", len(b))

	if err := sch.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
		return
	}
	t.Logf("%d docs loaded!", sch.DocCount())

	docs, infos = nil, nil
	sch.Search(nil, collector)
	t.Log("Docs:", docs, "Infos", infos)
	assert.Equal(t, "len(docs)()", len(docs), 1)
}

func BenchmarkTokenSetSearcher_1(b *testing.B) {
	log.Printf("1: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	ts1 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	for i := 0; i < 100000; i++ {
		if rand.Intn(2) == 0 {
			sch.AddDoc(ts1, i)
		} else {
			sch.AddDoc(ts0, i)
		}
	}
	//log.Printf("1: %d", b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sch.Search(ts1, func(int32, interface{}) error {
			return nil
		})
	}
}

func BenchmarkTokenSetSearcher_2(b *testing.B) {
	log.Printf("2: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]stringsp.Set{"text": stringsp.NewSet("A")}
	ts1 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	for i := 0; i < 100000; i++ {
		if rand.Intn(2) == 0 {
			sch.AddDoc(ts1, i)
		} else {
			sch.AddDoc(ts0, i)
		}
	}
	//log.Printf("2: %d", b.N)
	b.ResetTimer()
	N := b.N * 2
	for i := 0; i < N; i++ {
		sch.Search(ts1, func(int32, interface{}) error {
			return nil
		})
	}
}

func BenchmarkTokenSetSearcher_10(b *testing.B) {
	log.Printf("10: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]stringsp.Set{"text": stringsp.NewSet("A")}
	ts1 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	for i := 0; i < 100000; i++ {
		if rand.Intn(10) == 0 {
			sch.AddDoc(ts1, i)
		} else {
			sch.AddDoc(ts0, i)
		}
	}
	//log.Printf("10: %d", b.N)
	b.ResetTimer()
	N := b.N * 10
	for i := 0; i < N; i++ {
		sch.Search(ts1, func(int32, interface{}) error {
			return nil
		})
	}
}

func BenchmarkTokenSetSearcher_1000(b *testing.B) {
	log.Printf("1000: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]stringsp.Set{"text": stringsp.NewSet("A")}
	ts1 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	for i := 0; i < 100000; i++ {
		if rand.Intn(1000) == 0 {
			sch.AddDoc(ts1, i)
		} else {
			sch.AddDoc(ts0, i)
		}
	}
	//log.Printf("1000: %d", b.N)
	b.ResetTimer()
	N := b.N * 1000
	for i := 0; i < N; i++ {
		sch.Search(ts1, func(int32, interface{}) error {
			return nil
		})
	}
}

func BenchmarkTokenSetSearcher_100(b *testing.B) {
	log.Printf("100: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]stringsp.Set{"text": stringsp.NewSet("A")}
	ts1 := map[string]stringsp.Set{"text": stringsp.NewSet("A", "B")}
	for i := 0; i < 100000; i++ {
		if rand.Intn(100) == 0 {
			sch.AddDoc(ts1, i)
		} else {
			sch.AddDoc(ts0, i)
		}
	}
	//log.Printf("100: %d", b.N)
	b.ResetTimer()
	N := b.N * 100
	for i := 0; i < N; i++ {
		sch.Search(ts1, func(int32, interface{}) error {
			return nil
		})
	}
}

func TestTokenSetSearcher_bug1(t *testing.T) {
	DOCS := [][2]string{
		{" 0", "a b c"},
		{" 1", "a"},
		{" 2", "a"},
		{" 3", "a"},
		{" 4", "a b c"},
		{" 5", "a c"},
		{" 6", "a c"},
		{" 7", "a"},
		{" 8", "a c"},
	}
	sch := indexDocs(DOCS)

	var docs []int32
	var infos []*DocInfo
	collector := func(docID int32, data interface{}) error {
		docs = append(docs, docID)
		docInfo := data.(*DocInfo)
		infos = append(infos, docInfo)
		t.Logf("Doc: %d, %+v\n", docID, docInfo)
		return nil
	}

	sch.Search(SingleFieldQuery("text", "c", "b"), collector)
	assert.StringEqual(t, "docs", docs, "[0 4]")
}
