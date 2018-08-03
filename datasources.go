package iodata

import "github.com/moisespsena-go/iodata/api"

type DataSource struct {
	InputName  string
	DataHeader *api.DataHeader
}

func (i *DataSource) Name() string {
	return i.Name()
}

func (i *DataSource) Header() *api.DataHeader {
	return i.DataHeader
}

type Output struct {
	OutputName string
	DataHeader *api.DataHeader
}

func (i *Output) Name() string {
	return i.OutputName
}

func (i *Output) Header() *api.DataHeader {
	return i.DataHeader
}
