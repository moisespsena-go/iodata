package iodata

import (
	"bytes"

	"reflect"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

func ScanAssign(dataType api.DataType, data []byte, dest interface{}) (notBlank bool, err error) {
	destValue := reflect.ValueOf(dest).Elem()
	if bytes.Compare(data, dataType.BlankValue()) == 0 {
		if defaul := dataType.DefaultValue(); defaul != nil {
			destValue.Set(reflect.ValueOf(defaul))
		} else {
			destValue.Set(reflect.Zero(dataType.Type()))
		}
		return
	}
	notBlank = true
	var value interface{}
	if value, err = dataType.Scan(data); err == nil {
		destValue.Set(reflect.ValueOf(value))
	}
	return
}

func Scan(dataType api.DataType, data []byte) (v api.DataValue, err error) {
	ScanAssign(dataType, data, &v.Value)
	return
}

func ScanSliceAssign(dataType []api.DataType, data [][]byte, dest []interface{}) (notBlank []bool, err error) {
	notBlank = make([]bool, len(dest))
	for i, _ := range notBlank {
		notBlank[i], err = ScanAssign(dataType[i], data[i], dest[i])
		if err != nil {
			return nil, errwrap.Wrap(err, "Index %d", i)
		}
	}
	return
}

func ScanSliceStringsAssign(dataType []api.DataType, data []string, dest []interface{}) (notBlank []bool, err error) {
	notBlank = make([]bool, len(dest))
	for i, _ := range notBlank {
		notBlank[i], err = ScanAssign(dataType[i], []byte(data[i]), dest[i])
		if err != nil {
			return nil, errwrap.Wrap(err, "Index %d", i)
		}
	}
	return
}
