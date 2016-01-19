package index

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/golangplus/errors"
	"github.com/golangplus/testing/assert"
)

const N = 3

func createAndOpenBytesArr(t *testing.T, fn string) *ConstArrayReader {
	assert.NoErrorOrDie(t, os.RemoveAll(fn))
	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(t, err)

	for i := 0; i < N; i++ {
		s := fmt.Sprintf("data-%d", i)

		idx, err := w.AppendBytes([]byte(s))
		assert.NoErrorOrDie(t, err)

		assert.Equal(t, "idx", idx, i)
	}
	assert.NoErrorOrDie(t, w.Close())

	arr, err := OpenConstArray(fn)
	assert.NoErrorOrDie(t, err)
	assert.Equal(t, "len(arr.offsets)", len(arr.offsets), N+1)
	return arr
}

func TestConstArray_GetBytes(t *testing.T) {
	arr := createAndOpenBytesArr(t, path.Join(os.TempDir(), "./TestConstArray_ReadWrteBytes"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	for i := 0; i < N; i++ {
		t.Logf("i = %v", i)
		exp := fmt.Sprintf("data-%d", i)
		bs, err := arr.GetBytes(i)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, "s", string(bs), exp)
	}
}

func TestConstArray_FetchBytes(t *testing.T) {
	arr := createAndOpenBytesArr(t, path.Join(os.TempDir(), "./TestConstArray_FetchBytes"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	var indexes []int
	assert.NoError(t, arr.FetchBytes(func(index int, bs []byte) error {
		indexes = append(indexes, index)
		exp := fmt.Sprintf("data-%d", index)
		assert.Equal(t, "s", string(bs), exp)
		return nil
	}, 0, 2))
	assert.Equal(t, "indexes", indexes, []int{0, 2})

	// Check error returned.
	e := errors.New("inner-error")
	assert.Equal(t, "error", arr.FetchBytes(func(int, []byte) error {
		return e
	}, 0).(*errorsp.ErrorWithStacks).Err, e)
}

func TestConstArray_ForEachBytes(t *testing.T) {
	arr := createAndOpenBytesArr(t, path.Join(os.TempDir(), "./TestConstArray_ForEachBytes"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	var indexes []int
	assert.NoError(t, arr.ForEachBytes(func(index int, bs []byte) error {
		indexes = append(indexes, index)
		exp := fmt.Sprintf("data-%d", index)
		assert.Equal(t, "s", string(bs), exp)
		return nil
	}))
	assert.Equal(t, "indexes", indexes, []int{0, 1, 2})

	// Check error returned.
	e := errors.New("inner-error")
	assert.Equal(t, "error", arr.ForEachBytes(func(int, []byte) error {
		return e
	}).(*errorsp.ErrorWithStacks).Err, e)
}

func createAndOpenGobArr(t *testing.T, fn string) *ConstArrayReader {
	assert.NoErrorOrDie(t, os.RemoveAll(fn))
	w, err := CreateConstArray(fn)
	assert.NoErrorOrDie(t, err)

	for i := 0; i < N; i++ {
		s := fmt.Sprintf("data-%d", i)

		idx, err := w.AppendGob(s)
		assert.NoErrorOrDie(t, err)

		assert.Equal(t, "idx", idx, i)
	}
	assert.NoErrorOrDie(t, w.Close())

	arr, err := OpenConstArray(fn)
	assert.NoErrorOrDie(t, err)
	assert.Equal(t, "len(arr.offsets)", len(arr.offsets), N+1)
	return arr
}

func TestConstArray_GetGob(t *testing.T) {
	arr := createAndOpenGobArr(t, path.Join(os.TempDir(), "./TestConstArray_ReadWrteGob"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	for i := 0; i < N; i++ {
		t.Logf("i = %v", i)
		exp := fmt.Sprintf("data-%d", i)
		s, err := arr.GetGob(i)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, "s", s, exp)
	}
}

func TestConstArray_FetchGob(t *testing.T) {
	arr := createAndOpenGobArr(t, path.Join(os.TempDir(), "./TestConstArray_ReadWrteGob"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	var indexes []int
	assert.NoError(t, arr.FetchGobs(func(idx int, s interface{}) error {
		indexes = append(indexes, idx)
		exp := fmt.Sprintf("data-%d", idx)
		assert.Equal(t, "s", s, exp)
		return nil
	}, 0, 2))
	assert.Equal(t, "indexes", indexes, []int{0, 2})

	// Check error returned.
	e := errors.New("inner-error")
	assert.Equal(t, "error", arr.FetchGobs(func(int, interface{}) error {
		return e
	}, 0).(*errorsp.ErrorWithStacks).Err, e)
}

func TestConstArray_ForEachGob(t *testing.T) {
	arr := createAndOpenGobArr(t, path.Join(os.TempDir(), "./TestConstArray_ReadWrteGob"))
	defer func() {
		assert.NoError(t, arr.Close())
	}()

	var indexes []int
	assert.NoError(t, arr.ForEachGob(func(idx int, s interface{}) error {
		indexes = append(indexes, idx)
		exp := fmt.Sprintf("data-%d", idx)
		assert.Equal(t, "s", s, exp)
		return nil
	}))
	assert.Equal(t, "indexes", indexes, []int{0, 1, 2})

	// Check error returned.
	e := errors.New("inner-error")
	assert.Equal(t, "error", arr.ForEachGob(func(int, interface{}) error {
		return e
	}).(*errorsp.ErrorWithStacks).Err, e)
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
