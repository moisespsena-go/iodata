package csv

import (
	"encoding/csv"
	"io"

	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

type Reader struct {
	Reader *csv.Reader
}

func (r *Reader) Read(result [][][]byte) (count int, err error) {
	var record []string
	for count = range result {
		record, err = r.Reader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			err = errwrap.Wrap(err, "CSV Read")
			return
		}
		for i := range result[count] {
			result[count][i] = []byte(record[i])
		}
	}
	return
}

func NewReaderFactory() api.DataReaderFactory {
	return api.DataReaderFactoryFunc(func(header *api.DataHeader, args ...interface{}) api.DataReader {
		csvReader := csv.NewReader(args[0].(io.Reader))
		csvReader.Comma = ';'
		return &iodata.DataReader{BytesReader: &Reader{csvReader}, DataHeader: header}
	})
}

func NewWriterFactory() api.DataWriterFactory {
	return api.DataWriterFactoryFunc(func(header *api.DataHeader, args ...interface{}) api.DataWriter {
		csvWriter := csv.NewWriter(args[0].(io.Writer))
		csvWriter.Comma = ';'
		return &iodata.DataWriter{BytesWriter: &Writer{Writer: csvWriter}, DataHeader: header}
	})
}
