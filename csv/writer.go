package csv

import (
	"encoding/csv"

	"github.com/moisespsena-go/iodata"
)

type Writer struct {
	Writer        *csv.Writer
	writed        bool
	stringsWriter *iodata.StringsBytesWriter
}

func (w *Writer) Write(data [][][]byte) (err error) {
	if w.stringsWriter == nil {
		w.stringsWriter = &iodata.StringsBytesWriter{
			WriteFunc: func(strings [][]string) (err error) {
				return w.Writer.WriteAll(strings)
			},
		}
	}
	return w.stringsWriter.Write(data)
}
