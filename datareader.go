package iodata

import (
	"io"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

type DataReader struct {
	BytesReader api.BytesReader
	DataHeader  *api.DataHeader
	Count       int
	eof         bool
}

func (r *DataReader) Header() *api.DataHeader {
	return r.DataHeader
}

func (r *DataReader) Read(result ...[]interface{}) (count int, notBlank [][]bool, err error) {
	if r.eof {
		err = io.EOF
		return
	}

	l := len(result)
	byts := make([][][]byte, l)
	dataCount := len(result[0])
	for i := range result {
		byts[i] = make([][]byte, dataCount)
	}
	n := 0
	n, err = r.BytesReader.Read(byts)
	if n != l {
		byts = byts[0:n]
	}
	if err != nil {
		if err == io.EOF {
			r.eof = true
			return
		}
		err = errwrap.Wrap(err, "Read Bytes")
		return
	}

	typs := r.Header().Types
	notBlank = make([][]bool, l)
	for i := range byts {
		notBlank[i], err = ScanSliceAssign(typs, byts[i], result[i])
		if err != nil {
			err = errwrap.Wrap(err, "Scan Assign")
			return
		}
		r.Count++
		count++
	}
	return
}

func (r *DataReader) ReadOne(result ...interface{}) ([]bool, error) {
	_, notBlank, err := r.Read(result)
	if err != nil {
		return nil, err
	}
	return notBlank[0], nil
}

type DataReadCloser struct {
	api.DataReader
	CloseFunc func() error
}

func (r *DataReadCloser) Close() error {
	return r.CloseFunc()
}
