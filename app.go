package iodata

import "github.com/moisespsena-go/iodata/api"

type DataProcessorApplication struct {
	ReaderFactory api.DataReaderFactory
	WriterFactory api.DataWriterFactory
	Inputs        map[string]api.DataSource
}

func (app *DataProcessorApplication) Process(dp api.DataProcessorPlugin) error {
	outputHeader := dp.OutputHeader()
	dp.SetInputs(app.Inputs)
	w := app.WriterFactory.Factory(outputHeader)
	return dp.Process(w)
}
