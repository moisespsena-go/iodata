package db

import (
	"database/sql"

	"io"

	"reflect"

	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/modelstruct"
)

type Reader struct {
	DB         *sql.DB
	SQL        string
	SQLArgs    []interface{}
	DataHeader api.DataHeader
	open       bool
	Rows       *sql.Rows
	numRead    int
	eof        bool
}

func (r *Reader) Header() api.DataHeader {
	return r.DataHeader
}

func (r *Reader) NewValues(result []interface{}) (values []interface{}) {
	values = make([]interface{}, len(result))
	for i := range result {
		values[i] = NewRawScaner(&result[i]).SetToSlice(result, i)
	}
	return
}

func (r *Reader) preRead() (err error) {
	if r.eof {
		err = io.EOF
		return
	}
	if r.Rows == nil {
		r.Rows, err = r.DB.Query(r.SQL, r.SQLArgs...)
		if err != nil {
			return
		}
	}
	return
}

func (r *Reader) Read(result ...[]interface{}) (count int, notBlank [][]bool, err error) {
	return r.ReadF(r.NewValues, result...)
}

func (r *Reader) ReadF(newValues func(result []interface{}) (values []interface{}), result ...[]interface{}) (count int, notBlank [][]bool, err error) {
	if err = r.preRead(); err != nil {
		return
	}

	for _, res := range result {
		if !r.Rows.Next() {
			r.eof = true
			return count, notBlank, io.EOF
		}
		values := newValues(res)
		err = r.Rows.Scan(values...)
		if err != nil {
			return
		}
		count++
	}
	return
}

func (r *Reader) ReadModelStruct(model *modelstruct.ModelStruct, results ...interface{}) (count int, err error) {
	if r.Rows == nil {
		if err = r.preRead(); err != nil {
			return
		}
	}

	var (
		resultColumns []string
		modelValue    reflect.Value
	)

	if resultColumns, err = r.Rows.Columns(); err != nil {
		return
	}

	fieldScaner := func(modelValue reflect.Value, i int, field *modelstruct.StructField) interface{} {
		value := modelValue.FieldByIndex(field.Index)
		rs := NewRawScaner(value.Addr().Interface())
		rs.Set = func(rs *RawScaner) {
			data := reflect.ValueOf(rs.Data)
			value.Set(data.Elem())
		}
		return rs
	}

	var dbResults = make([][]interface{}, len(results))

	for i := range results {
		if results[i] == nil {
			modelValue = reflect.New(model.Type).Elem()
			results[i] = modelValue.Addr().Interface()
		} else {
			modelValue = reflect.ValueOf(results[i]).Elem()
		}

		dbResults[i] = make([]interface{}, len(resultColumns))

		for j, column := range resultColumns {
			if field, ok := model.FieldsByName[column]; ok {
				dbResults[i][j] = fieldScaner(modelValue, j, field)
			} else {
				dbResults[i][j] = &RawScanerDiscard{}
			}
		}
	}

	count, _, err = r.ReadF(func(result []interface{}) (values []interface{}) {
		return result
	}, dbResults...)
	return
}

func (r *Reader) ReadStruct(results ...interface{}) (count int, err error) {
	return iodata.ReadStruct(r.ReadModelStruct, results...)
}

func (r *Reader) ReadOne(result ...interface{}) ([]bool, error) {
	_, _, err := r.Read(result)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Reader) Close() error {
	defer func() {
		r.Rows = nil
	}()
	if r.Rows != nil {
		return r.Rows.Close()
	}
	return nil
}
