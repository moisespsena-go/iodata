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
	searcher      api.Searcher
}

func (ds *DataSource) Searcher() api.Searcher {
	panic("Not Implemented")
}

func (i *DataSource) Load() (api.DataReadCloser, error) {
	r, err := os.Open(i.FilePath)
	if err != nil {
		return nil, err
	}
	dr := i.ReaderFactory.Factory(i.DataHeader, r)
	return &iodata.DataReadCloser{dr, r.Close}, nil
}