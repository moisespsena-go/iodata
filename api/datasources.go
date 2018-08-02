package api

type DataSource interface {
	Name() string
	Filter(query string, args ...interface{})
	Load() (DataReadCloser, error)
	Header() *DataHeader
}

type Output interface {
	Name() string
	Writer() (DataWriteCloser, error)
	Header() *DataHeader
}