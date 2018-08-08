package db

import (
	"database/sql"

	"io"

	"github.com/moisespsena-go/iodata/api"
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

func (r *Reader) Read(result ...[]interface{}) (count int, notBlank [][]bool, err error) {
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
	if len(result) == 0 {
		return
	}

	for _, res := range result {
		if !r.Rows.Next() {
			r.eof = true
			return count, notBlank, io.EOF
		}
		values := r.NewValues(res)
		err = r.Rows.Scan(values...)
		if err != nil {
			return
		}
		count++
	}
	return
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
