package iodata

import (
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

type DataWriter struct {
	BytesWriter api.BytesWriter
	DataHeader  api.DataHeader
	Count       int
}

func (r *DataWriter) Header() api.DataHeader {
	return r.DataHeader
}

func (w *DataWriter) Write(data ...[]interface{}) (err error) {
	byts, err := Dump(w.Header(), data)
	if err != nil {
		return errwrap.Wrap(err, "Dump")
	}
	return w.BytesWriter.Write(byts)
}

func (w *DataWriter) WriteP(data ...[]interface{}) (err error) {
	byts, err := DumpP(w.Header(), data)
	if err != nil {
		return errwrap.Wrap(err, "DumpP")
	}
	return w.BytesWriter.Write(byts)
}

func (w *DataWriter) WriteOne(data ...interface{}) (err error) {
	return w.Write(data)
}

type DataWriteCloser struct {
	api.DataWriter
	CloseFunc func() error
}

func (r *DataWriteCloser) Close() error {
	return r.CloseFunc()
}
