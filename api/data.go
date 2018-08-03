package api

import "io"

type DataHeader struct {
	Names []string
	Types []DataType
	ByName map[string]int
}

func NewDataHeader(names []string, types []DataType) *DataHeader {
	dh := &DataHeader{names, types, map[string]int{}}
	for i, name := range names {
		dh.ByName[name] = i
	}
	return dh
}

type DataReader interface {
	Header() *DataHeader
	Read(result ...[]interface{}) (count int, notBlank [][]bool, err error)
	ReadOne(result ...interface{}) (notBlank []bool, err error)
}

type DataReadCloser interface {
	DataReader
	io.Closer
}

type DataWriter interface {
	Header() *DataHeader
	Write(data ...[]interface{}) (err error)
	WriteOne(data ...interface{}) (err error)
}

type DataWriteCloser interface {
	DataWriter
	io.Closer
}