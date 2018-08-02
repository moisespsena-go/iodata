package api

import "io"

type DataHeader struct {
	Names []string
	Types []DataType
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