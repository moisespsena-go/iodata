package csv

import (
	"os"

	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api"
)

type DataSource struct {
	iodata.DataSource
	ReaderFactory api.DataReaderFactory
	FilePath      string
}

func (i *DataSource) Load() (api.DataReadCloser, error) {
	r, err := os.Open(i.FilePath)
	if err != nil {
		return nil, err
	}
	dr := i.ReaderFactory.Factory(i.DataHeader, r, i.Filters)
	return &iodata.DataReadCloser{dr, r.Close}, nil
}

type Output struct {
	iodata.Output
	WriterFactory api.DataWriterFactory
	FilePath      string
}

func (i *Output) Writer() (api.DataWriteCloser, error) {
	w, err := os.Create(i.FilePath)
	if err != nil {
		return nil, err
	}
	dw := i.WriterFactory.Factory(i.DataHeader, w)
	return &iodata.DataWriteCloser{dw, w.Close}, nil
}
