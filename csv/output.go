package csv

import (
	"os"

	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api"
)

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
