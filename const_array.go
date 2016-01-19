package index

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
	"path"

	"github.com/golangplus/bytes"
	"github.com/golangplus/errors"
)

type ConstArrayReader struct {
	offsets   []int64
	dataFiles chan *os.File
}

func OpenConstArray(dir string) (*ConstArrayReader, error) {
	of, err := os.Open(path.Join(dir, saOffsetsFilename))
	if err != nil {
		return nil, errorsp.WithStacks(err)
	}
	defer of.Close()
	l, err := of.Seek(0, 2)
	if err != nil {
		return nil, errorsp.WithStacks(err)
	}
	offsets := make([]int64, l/8)
	if _, err := of.Seek(0, 0); err != nil {
		return nil, errorsp.WithStacks(err)
	}
	for i := range offsets {
		if err := binary.Read(of, binary.BigEndian, &offsets[i]); err != nil {
			if err == io.EOF {
				if i == len(offsets)-1 {
					break
				}
				return nil, errorsp.WithStacks(io.ErrUnexpectedEOF)
			}
			return nil, errorsp.WithStacks(err)
		}
	}
	dfs := make(chan *os.File, 10)
	for i := 0; i < cap(dfs); i++ {
		df, err := os.Open(path.Join(dir, saDataFilename))
		if err != nil {
			return nil, errorsp.WithStacks(err)
		}
		dfs <- df
	}
	w := &ConstArrayReader{
		offsets:   offsets,
		dataFiles: dfs,
	}
	return w, nil
}

func (r *ConstArrayReader) Close() error {
	var err error
	for i := 0; i < cap(r.dataFiles); i++ {
		df := <-r.dataFiles
		if e := df.Close(); e != nil {
			err = e
		}
	}
	return err
}

func (r *ConstArrayReader) returnDataFile(df *os.File) {
	r.dataFiles <- df
}

func (r *ConstArrayReader) GetBytes(index int) ([]byte, error) {
	df := <-r.dataFiles
	if df == nil {
		return nil, nil
	}
	defer r.returnDataFile(df)

	_, err := df.Seek(r.offsets[index], 0)
	if err != nil {
		return nil, errorsp.WithStacks(err)
	}
	bs := make([]byte, r.offsets[index+1]-r.offsets[index])
	if _, err := io.ReadFull(df, bs); err != nil {
		return nil, errorsp.WithStacks(err)
	}
	return bs, nil
}

func (r *ConstArrayReader) FetchBytes(output func(int, []byte) error, indexes ...int) error {
	df := <-r.dataFiles
	if df == nil {
		return nil
	}
	defer r.returnDataFile(df)

	for _, index := range indexes {
		_, err := df.Seek(r.offsets[index], 0)
		if err != nil {
			return errorsp.WithStacks(err)
		}
		bs := make([]byte, r.offsets[index+1]-r.offsets[index])
		if _, err := io.ReadFull(df, bs); err != nil {
			return errorsp.WithStacks(err)
		}
		if err := output(index, bs); err != nil {
			return errorsp.WithStacks(err)
		}
	}
	return nil
}

func (r *ConstArrayReader) ForEachBytes(output func(int, []byte) error) error {
	df := <-r.dataFiles
	if df == nil {
		return nil
	}
	defer r.returnDataFile(df)

	if _, err := df.Seek(0, 0); err != nil {
		return errorsp.WithStacks(err)
	}
	for i := 1; i < len(r.offsets); i++ {
		bs := make([]byte, r.offsets[i]-r.offsets[i-1])
		n, err := df.Read(bs)
		if err != nil || n != len(bs) {
			if err == io.EOF || n != len(bs) {
				return errorsp.WithStacks(io.ErrUnexpectedEOF)
			}
			return errorsp.WithStacks(err)
		}
		if err := output(i-1, bs); err != nil {
			return errorsp.WithStacks(err)
		}
	}
	return nil
}

func (r *ConstArrayReader) GetGob(index int) (interface{}, error) {
	df := <-r.dataFiles
	if df == nil {
		return nil, nil
	}
	defer r.returnDataFile(df)

	if _, err := df.Seek(r.offsets[index], 0); err != nil {
		return nil, errorsp.WithStacks(err)
	}
	var e interface{}
	if err := gob.NewDecoder(df).Decode(&e); err != nil {
		return nil, errorsp.WithStacks(err)
	}
	return e, nil
}

func (r *ConstArrayReader) FetchGobs(output func(int, interface{}) error, indexes ...int) error {
	df := <-r.dataFiles
	if df == nil {
		return nil
	}
	defer r.returnDataFile(df)

	for _, index := range indexes {
		if _, err := df.Seek(r.offsets[index], 0); err != nil {
			return errorsp.WithStacks(err)
		}
		var e interface{}
		if err := gob.NewDecoder(df).Decode(&e); err != nil {
			return errorsp.WithStacks(err)
		}
		if err := output(index, e); err != nil {
			return errorsp.WithStacks(err)
		}
	}
	return nil
}

func (r *ConstArrayReader) ForEachGob(output func(int, interface{}) error) error {
	df := <-r.dataFiles
	if df == nil {
		return nil
	}
	defer r.returnDataFile(df)

	if _, err := df.Seek(0, 0); err != nil {
		return errorsp.WithStacks(err)
	}
	dec := gob.NewDecoder(df)
	for i := 1; i < len(r.offsets); i++ {
		var e interface{}
		if err := dec.Decode(&e); err != nil {
			return errorsp.WithStacks(err)
		}
		if err := output(i-1, e); err != nil {
			return errorsp.WithStacks(err)
		}
	}
	return nil
}

type ConstArrayWriter struct {
	count    int
	offset   int64
	offsFile *os.File
	dataFile *os.File
}

func CreateConstArray(dir string) (*ConstArrayWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errorsp.WithStacks(err)
	}
	var of, df *os.File
	defer func() {
		if of != nil {
			of.Close()
		}
		if df != nil {
			df.Close()
		}
	}()
	of, err := os.Create(path.Join(dir, saOffsetsFilename))
	if err != nil {
		return nil, errorsp.WithStacks(err)
	}
	df, err = os.Create(path.Join(dir, saDataFilename))
	if err != nil {
		return nil, errorsp.WithStacks(err)
	}
	w := &ConstArrayWriter{
		count:    0,
		offset:   0,
		offsFile: of,
		dataFile: df,
	}
	of, df = nil, nil
	return w, nil
}

const (
	saDataFilename    = "data"
	saOffsetsFilename = "offsets"
)

func (sa *ConstArrayWriter) Close() error {
	var err error
	if e := binary.Write(sa.offsFile, binary.BigEndian, sa.offset); e != nil {
		err = e
	}
	if e := sa.offsFile.Close(); e != nil {
		err = e
	}
	if e := sa.dataFile.Close(); e != nil {
		err = e
	}
	return errorsp.WithStacks(err)
}

func (sa *ConstArrayWriter) AppendBytes(bs []byte) (int, error) {
	if err := binary.Write(sa.offsFile, binary.BigEndian, sa.offset); err != nil {
		return 0, errorsp.WithStacks(err)
	}
	_, err := sa.dataFile.Write(bs)
	if err != nil {
		return 0, errorsp.WithStacks(err)
	}
	sa.count++
	sa.offset += int64(len(bs))
	return sa.count - 1, nil
}

func (sa *ConstArrayWriter) AppendGob(e interface{}) (int, error) {
	var bs bytesp.Slice
	if err := gob.NewEncoder(&bs).Encode(&e); err != nil {
		return 0, errorsp.WithStacks(err)
	}
	return sa.AppendBytes(bs)
}
