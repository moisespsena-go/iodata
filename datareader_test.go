package iodata

import (
	"testing"

	"io"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/api/datatypes"
)

func TestDataReader_Read(t *testing.T) {
	var (
		br       = &StringsBytesReader{Reader: NewSliceStringsReader([]string{"15", "1.2"}, []string{"7", "5.0896"})}
		types    = []api.DataType{&datatypes.Int64Type{}, &datatypes.Float64Type{}}
		r        = &DataReader{BytesReader: br, DataHeader: &api.DataHeader{Types: types}}
		a0, a1   int64
		b0, b1   float64
		result         = [][]interface{}{{&a0, &b0}, {&a1, &b1}}
		ea0, ea1 int64 = 15, 7
		eb0, eb1       = 1.2, 5.0896
	)
	count, notBlank, err := r.Read(result...)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("count != 2")
	}
	if len(notBlank) != 2 {
		t.Errorf("len(notBlank) != 2")
	}
	if len(notBlank[0]) != 2 {
		t.Errorf("len(notBlank[i]) != 2")
	}
	if !notBlank[0][0] {
		t.Errorf("notBlank[0][0] is blank")
	}
	if !notBlank[0][1] {
		t.Errorf("notBlank[0][1] is blank")
	}
	if a0 != ea0 {
		t.Errorf("a0 != %v", ea0)
	}
	if b0 != eb0 {
		t.Errorf("b0 != %v", eb0)
	}
	if a1 != ea1 {
		t.Errorf("a1 != %v", ea1)
	}
	if b1 != eb1 {
		t.Errorf("b1 != %v", eb1)
	}
	count, notBlank, err = r.Read(result...)
	if err == nil {
		t.Errorf("Error Expected")
	}
	if err != io.EOF {
		t.Errorf("EOF Expected: %v", err)
	}
	if count != 0 {
		t.Errorf("Count != 0 expected")
	}
	if len(notBlank) != 0 {
		t.Errorf("len(notBlank) != 0")
	}
}

func TestDataReader_Read2(t *testing.T) {
	var (
		br    = &StringsBytesReader{Reader: NewSliceStringsReader([]string{"15", "1.2"}, []string{"7", "5.0896"})}
		types = []api.DataType{&datatypes.Int64Type{}, &datatypes.Float64Type{}}
		r     = &DataReader{BytesReader: br, DataHeader: &api.DataHeader{Types: types}}
		a0    int64
		b0    float64
	)
	count, notBlank, err := r.Read([]interface{}{&a0, &b0})
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("count != 1")
	}
	if len(notBlank) != 1 {
		t.Errorf("len(notBlank) != 1")
	}
	count, notBlank, err = r.Read([]interface{}{&a0, &b0})
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("count != 1")
	}
	if len(notBlank) != 1 {
		t.Errorf("len(notBlank) != 1")
	}
	count, notBlank, err = r.Read([]interface{}{&a0, &b0})
	if err == nil {
		t.Errorf("Error Expected")
	}
	if err != io.EOF {
		t.Errorf("EOF Expected: %v", err)
	}
	if count != 0 {
		t.Errorf("count != 0")
	}
	if len(notBlank) != 0 {
		t.Errorf("len(notBlank) != 0")
	}
}
