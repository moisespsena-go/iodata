package api

type DataSource interface {
	Name() string
	Load() (DataReadCloser, error)
	Header() DataHeader
	Searcher()Searcher
}

type Output interface {
	Name() string
	Writer() (DataWriteCloser, error)
	Header() DataHeader
}