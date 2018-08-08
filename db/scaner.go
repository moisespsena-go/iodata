package db

import (
	"database/sql"
	"reflect"
)

type RawScaner struct {
	Typ     reflect.Type
	Data    interface{}
	NotNil  bool
	IsValid bool
	IsPtr   bool
	Set     func(f *RawScaner)
}

func NewRawScaner(defaul interface{}, set ...func(rs *RawScaner)) *RawScaner {
	rs := &RawScaner{Typ: reflect.TypeOf(defaul).Elem(), Data: defaul}
	if len(set) > 0 {
		rs.Set = set[0]
	}
	return rs
}

func (f *RawScaner) IsNil() bool {
	return !f.NotNil
}

func (f *RawScaner) Scan(src interface{}) error {
	if f.NotNil = src != nil; f.NotNil {
		f.Data = reflect.New(f.Typ).Interface()
		if scan, ok := f.Data.(sql.Scanner); ok {
			return scan.Scan(src)
		}
		err := convertAssign(f.Data, src)
		if err == nil && f.Set != nil {
			f.Set(f)
		}
		return err
	}
	return nil
}

func (rs *RawScaner) SetToSlice(dest []interface{}, i int) *RawScaner {
	rs.Set = func(rs *RawScaner) {
		res := reflect.ValueOf(dest[i]).Elem()
		data := reflect.ValueOf(rs.Data).Elem()
		res.Set(data.Elem())
	}
	return rs
}
