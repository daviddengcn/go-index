package index

import (
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-villa"
	"log"
	"math/rand"
	"testing"
)

type DocInfo struct {
	A string
}

func TestTokenSetSearcher(t *testing.T) {
	DOCS := [][]string{
		{"To friends", "hello my friend"},
		{"To dogs", "GO go go, my dog"},
	}

	sch := &TokenSetSearcher{}

	for i := range DOCS {
		text := DOCS[i][1]
		var tokens villa.StrSet
		TokenizeBySeparators(" ,", villa.NewPByteSlice([]byte(text)),
			func(token []byte) error {
				tokens.Put(string(token))
				return nil
			})

		fields := map[string]villa.StrSet{
			"text": tokens,
		}
		sch.AddDoc(fields, &DocInfo{
			A: fmt.Sprintf("%d - %s", i+1, DOCS[i][0]),
		})
	}

	var docs []int32
	var infos []*DocInfo
	collector := func(docID int32, data interface{}) error {
		docs = append(docs, docID)
		docInfo := data.(*DocInfo)
		infos = append(infos, docInfo)
		fmt.Printf("Doc: %d, %+v\n", docID, docInfo)
		return nil
	}
	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "my"), collector)
	fmt.Println("Docs:", docs, "Infos", infos)
	villa.AssertEquals(t, "len(docs)(my)", len(docs), 2)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "my", "dog"), collector)
	fmt.Println("Docs:", docs, "Infos", infos)
	villa.AssertEquals(t, "len(docs)(my dog)", len(docs), 1)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "friend"), collector)
	fmt.Println("Docs:", docs, "Infos", infos)
	villa.AssertEquals(t, "len(docs)(friend)", len(docs), 1)

	sch.Delete(0)

	docs, infos = nil, nil
	sch.Search(nil, collector)
	fmt.Println("Docs:", docs, "Infos", infos)
	villa.AssertEquals(t, "len(docs)()", len(docs), 1)

	var b villa.ByteSlice
	gob.Register(&DocInfo{})
	if err := sch.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
	}
	log.Printf("%d bytes written", len(b))

	if err := sch.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
	}
	log.Printf("%d docs loaded!", sch.DocCount())

	docs, infos = nil, nil
	sch.Search(nil, collector)
	fmt.Println("Docs:", docs, "Infos", infos)
	villa.AssertEquals(t, "len(docs)()", len(docs), 1)
}

func BenchmarkTokenSetSearcher_1(b *testing.B) {
	log.Printf("1: %d", b.N)
	sch := &TokenSetSearcher{}
	rand.Seed(1)
	ts0 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
	ts1 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
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
	ts0 := map[string]villa.StrSet{"text": villa.NewStrSet("A")}
	ts1 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
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
	ts0 := map[string]villa.StrSet{"text": villa.NewStrSet("A")}
	ts1 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
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
	ts0 := map[string]villa.StrSet{"text": villa.NewStrSet("A")}
	ts1 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
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
	ts0 := map[string]villa.StrSet{"text": villa.NewStrSet("A")}
	ts1 := map[string]villa.StrSet{"text": villa.NewStrSet("A", "B")}
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
