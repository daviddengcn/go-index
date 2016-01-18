package index

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golangplus/testing/assert"
)

func TestConstArray_ReadWrteBytes(t *testing.T) {
	fn := "./TestConstArray_ReadWrteGob"
	assert.NoErrorOrDie(t, os.RemoveAll(fn))

	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(t, err)
	const N = 100
	for i := 0; i < N; i++ {
		s := fmt.Sprintf("data-%d", i)
		idx, err := w.AppendBytes([]byte(s))
		assert.NoError(t, err)
		assert.Equal(t, "idx", idx, i)
	}
	assert.NoError(t, w.Close())

	arr, err := OpenConstArray(fn)
	assert.NoErrorOrDie(t, err)
	assert.Equal(t, "len(arr.offsets)", len(arr.offsets), N+1)
	for i := 0; i < N; i++ {
		t.Logf("i = %v", i)
		exp := fmt.Sprintf("data-%d", i)
		bs, err := arr.GetBytes(i)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, "s", string(bs), exp)
	}
	assert.NoError(t, arr.Close())
}

func TestConstArray_ReadWrteGob(t *testing.T) {
	fn := "./TestConstArray_ReadWrteGob"
	assert.NoErrorOrDie(t, os.RemoveAll(fn))

	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(t, err)
	const N = 100
	for i := 0; i < N; i++ {
		s := fmt.Sprintf("data-%d", i)
		idx, err := w.AppendGob(s)
		assert.NoError(t, err)
		assert.Equal(t, "idx", idx, i)
	}
	assert.NoError(t, w.Close())

	arr, err := OpenConstArray(fn)
	assert.NoErrorOrDie(t, err)
	assert.Equal(t, "len(arr.offsets)", len(arr.offsets), N+1)
	for i := 0; i < N; i++ {
		t.Logf("i = %v", i)
		exp := fmt.Sprintf("data-%d", i)
		s, err := arr.GetGob(i)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, "s", s, exp)
	}

	var idxs []int
	var daList []interface{}
	assert.NoError(t, arr.ForEachGob(func(idx int, data interface{}) error {
		idxs = append(idxs, idx)
		daList = append(daList, data)
		return nil
	}))
	for i := 0; i < N; i++ {
		exp := fmt.Sprintf("data-%d", i)
		assert.Equal(t, "idxs[i]", idxs[i], i)
		assert.Equal(t, "daList[i]", daList[i], exp)
	}

	assert.NoError(t, arr.Close())
}

func BenchmarkConstArrayIndexing(b *testing.B) {
	fn := "./BenchmarkTestConstArrayIndexing"
	assert.NoErrorOrDie(b, os.RemoveAll(fn))

	log.Printf("1: %d", b.N)
	b.ResetTimer()
	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(b, err)
	for i := 0; i < b.N; i++ {
		w.AppendGob(i)
	}
}

func BenchmarkConstArrayRead(b *testing.B) {
	fn := "./BenchmarkTestConstArrayIndexing"
	assert.NoErrorOrDie(b, os.RemoveAll(fn))

	log.Printf("1: %d", b.N)
	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(b, err)
	for i := 0; i < b.N; i++ {
		w.AppendGob(i)
	}
	w.Close()
	r, err := OpenConstArray(fn)
	assert.NoErrorOrDie(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.GetGob(i)
	}
	b.StopTimer()
	assert.NoError(b, r.Close())
}
