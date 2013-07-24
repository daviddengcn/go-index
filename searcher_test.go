package index

import (
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-algs/ed"
	"github.com/daviddengcn/go-villa"
	"log"
	"math/rand"
	"strings"
	"testing"
)

type DocInfo struct {
	A string
}

func indexDocs(docs [][2]string) *TokenSetSearcher {
	sch := &TokenSetSearcher{}
	for i := range docs {
		text := docs[i][1]
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
	AssertEquals(t, "len(docs)(my)", len(docs), 2)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "my", "dog"), collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	AssertEquals(t, "len(docs)(my dog)", len(docs), 1)

	docs, infos = nil, nil
	sch.Search(SingleFieldQuery("text", "friend"), collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	AssertEquals(t, "len(docs)(friend)", len(docs), 1)

	docs = sch.TokenDocList("text", "my")
	AssertStringEquals(t, "text:To", docs, "[0 1]")

	sch.Delete(0)

	docs, infos = nil, nil
	sch.Search(nil, collector)
	//fmt.Println("Docs:", docs, "Infos", infos)
	AssertEquals(t, "len(docs)()", len(docs), 1)

	var b villa.ByteSlice
	gob.Register(&DocInfo{})
	if err := sch.Save(&b); err != nil {
		t.Errorf("Save failed: %v", err)
	}
	t.Logf("%d bytes written", len(b))

	if err := sch.Load(&b); err != nil {
		t.Errorf("Load failed: %v", err)
	}
	t.Logf("%d docs loaded!", sch.DocCount())

	docs, infos = nil, nil
	sch.Search(nil, collector)
	t.Logf("Docs:", docs, "Infos", infos)
	AssertEquals(t, "len(docs)()", len(docs), 1)
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
	AssertStringEquals(t, "docs", docs, "[0 4]")
}

func AssertTextEquals(t *testing.T, name, act, exp string) {
	if exp == act {
		return
	}

	expLines := strings.Split(exp, "\n")
	actLines := strings.Split(act, "\n")

	t.Errorf("%s unexpected(exp: %d lines, act %d lines!", name, len(expLines), len(actLines))
	t.Logf("exp ---  act +++")
	t.Logf("Difference:")
	_, matA, matB := ed.EditDistanceFFull(len(expLines), len(actLines), func(iA, iB int) int {
		sa, sb := expLines[iA], actLines[iB]
		if sa == sb {
			return 0
		}
		return ed.String(sa, sb)
	}, func(iA int) int {
		return len(expLines[iA]) + 1
	}, func(iB int) int {
		return len(actLines[iB]) + 1
	})
	for i, j := 0, 0; i < len(expLines) || j < len(actLines); {
		switch {
		case j >= len(actLines) || i < len(expLines) && matA[i] < 0:
			t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
			i++
		case i >= len(expLines) || j < len(actLines) && matB[j] < 0:
			t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			j++
		default:
			if expLines[i] != actLines[j] {
				t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
				t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			} // else
			i++
			j++
		}
	} // for i, j
}

func AssertStringsEqual(t *testing.T, name string, act, exp []string) {
	if villa.StringSlice(exp).Equals(act) {
		return
	}
	t.Errorf("%s unexpected(exp: %d lines, act %d lines)!", name, len(exp), len(act))
	t.Logf("exp ---  act +++")
	t.Logf("Difference:")
	_, matA, matB := ed.EditDistanceFFull(len(exp), len(act), func(iA, iB int) int {
		sa, sb := exp[iA], act[iB]
		if sa == sb {
			return 0
		}
		return ed.String(sa, sb)
	}, func(iA int) int {
		return len(exp[iA]) + 1
	}, func(iB int) int {
		return len(act[iB]) + 1
	})
	for i, j := 0, 0; i < len(exp) || j < len(act); {
		switch {
		case j >= len(act) || i < len(exp) && matA[i] < 0:
			t.Logf("--- %3d: %s", i+1, showText(exp[i]))
			i++
		case i >= len(exp) || j < len(act) && matB[j] < 0:
			t.Logf("+++ %3d: %s", j+1, showText(act[j]))
			j++
		default:
			if exp[i] != act[j] {
				t.Logf("--- %3d: %s", i+1, showText(exp[i]))
				t.Logf("+++ %3d: %s", j+1, showText(act[j]))
			} // else
			i++
			j++
		}
	} // for i, j
}

