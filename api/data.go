package api

import "io"

type DataField interface {
	Name() string
	Type() DataType
}

type dataField struct {
	name string
	typ  DataType
}

func (d *dataField) Name() string {
	return d.name
}

func (d *dataField) Type() DataType {
	return d.typ
}

func NewDataField(name string, typ DataType) DataField {
	return &dataField{name, typ}
}

type DataHeader interface {
	Fields() []DataField
	IndexOf(fieldName string) int
	HasField(fieldName string) bool
	Field(fieldName string) DataField
	Names() []string
	Types() []DataType
	TypeNames() []DataTypeName
}

type dataHeader struct {
	fields []DataField
	index  map[string]int
}

func (dh dataHeader) Fields() []DataField {
	return dh.fields
}

func (dh dataHeader) IndexOf(fieldName string) int {
	if i, ok := dh.index[fieldName]; ok {
		return i
	}
	return -1
}

func (dh dataHeader) HasField(fieldName string) bool {
	if _, ok := dh.index[fieldName]; ok {
		return true
	}
	return false
}

func (dh dataHeader) Field(fieldName string) DataField {
	if i, ok := dh.index[fieldName]; ok {
		return dh.fields[i]
	}
	return nil
}

func (dh dataHeader) Names() []string {
	names := make([]string, len(dh.fields))
	for i, f := range dh.fields {
		names[i] = f.Name()
	}
	return names
}

func (dh dataHeader) Types() []DataType {
	types := make([]DataType, len(dh.fields))
	for i, f := range dh.fields {
		types[i] = f.Type()
	}
	return types
}

func (dh dataHeader) TypeNames() []DataTypeName {
	names := make([]DataTypeName, len(dh.fields))
	for i, f := range dh.fields {
		names[i] = f.Type().Name()
	}
	return names
}

func NewDataHeader(fields ...DataField) DataHeader {
	dh := &dataHeader{fields, map[string]int{}}
	for i, f := range fields {
		dh.index[f.Name()] = i
	}
	return dh
}

type DataReader interface {
	Header() DataHeader
	Read(result ...[]interface{}) (count int, notBlank [][]bool, err error)
	ReadOne(result ...interface{}) (notBlank []bool, err error)
}

type DataReadCloser interface {
	DataReader
	io.Closer
}

type DataWriter interface {
	Header() DataHeader
	Write(data ...[]interface{}) (err error)
	WriteOne(data ...interface{}) (err error)
}

type DataWriteCloser interface {
	DataWriter
	io.Closer
}
