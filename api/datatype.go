package api

import (
	"fmt"
	"reflect"
)

type DataTypeName string

func (dt DataTypeName) IsSlice() bool {
	return dt[0:2] == "[]"
}

func (dt DataTypeName) Elem() DataTypeName {
	if dt.IsSlice() {
		return dt[2:]
	}
	panic(fmt.Errorf("DataType %q isn't slice.", dt))
}

type DataType interface {
	Name() DataTypeName
	IsSlice() bool
	Scan(value []byte) (interface{}, error)
	Dump(value interface{}) ([]byte, error)
	Elem() DataType
	DefaultValue() interface{}
	BlankValue() []byte
	Type() reflect.Type
}

type DataTypeBase struct {
	TypeName DataTypeName
	Default  interface{}
	Typ      reflect.Type
}

func (b *DataTypeBase) Name() DataTypeName {
	return b.TypeName
}

func (b *DataTypeBase) IsSlice() bool {
	return b.TypeName.IsSlice()
}

func (b *DataTypeBase) BlankValue() (v []byte) {
	return v
}

func (d *DataTypeBase) DefaultValue() interface{} {
	return d.Default
}

func (d *DataTypeBase) Type() reflect.Type {
	return d.Typ
}
