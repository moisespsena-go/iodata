package api

type DataReaderFactory interface {
	Factory(header DataHeader, args ...interface{}) DataReader
}

type DataWriterFactory interface {
	Factory(header DataHeader, args ...interface{}) DataWriter
}

type DataReaderFactoryFunc func(header DataHeader, args ...interface{}) DataReader

func (f DataReaderFactoryFunc) Factory(header DataHeader, args ...interface{}) DataReader {
	return f(header, args...)
}

type DataWriterFactoryFunc func(header DataHeader, args ...interface{}) DataWriter

func (f DataWriterFactoryFunc) Factory(header DataHeader, args ...interface{}) DataWriter {
	return f(header, args...)
}

type StringsReader interface {
	Read([]string) error
}

type BytesReader interface {
	Read(result [][][]byte) (count int, err error)
}

type BytesWriter interface {
	Write(data [][][]byte) (err error)
}

type StringsWriter interface {
	Write([][]string) error
}
