package iodata

import (
	"io"

	"reflect"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/modelstruct"
	"github.com/moisespsena/go-error-wrap"
)

type DataReader struct {
	BytesReader api.BytesReader
	DataHeader  api.DataHeader
	Count       int
	eof         bool
}

func (r *DataReader) Header() api.DataHeader {
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

	typs := r.Header().Types()
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

func ReadStruct(reader func(model *modelstruct.ModelStruct, results ...interface{}) (count int, err error), results ...interface{}) (count int, err error) {
	if len(results) == 0 {
		return
	}

	rt := reflect.TypeOf(results[0])
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.Slice:
		// get type of ptr value at slice[0]: []*User{&u} or make([]*User, n)
		model := modelstruct.Get(rt.Elem().Elem())
		var (
			rv              = reflect.ValueOf(results[0])
			resultsInteface = make([]interface{}, rv.Len())
			elem            reflect.Value
		)

		for i := range resultsInteface {
			if elem = rv.Index(i); elem.IsNil() {
				elem.Set(reflect.New(model.Type))
			}
			resultsInteface[i] = elem.Interface()
		}

		return reader(model, resultsInteface...)
	case reflect.Struct:
		model := modelstruct.Get(rt)
		return reader(model, results...)
	}
	return
}
