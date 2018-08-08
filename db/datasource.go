package db

import (
	"bytes"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/query"
)

type DataSource struct {
	iodata.DataSource
	ReaderFactory api.DataReaderFactory
	searcher      *iodata.Search
	db            *sql.DB
}

func NewDataSource(db *sql.DB) *DataSource {
	return &DataSource{db: db}
}

func (ds *DataSource) Searcher() api.Searcher {
	if ds.searcher == nil {
		ds.searcher = &iodata.Search{}
	}
	return ds.searcher
}

func (i *DataSource) NewSearchScope() (*iodata.Scope, error) {
	dialect, _ := gorm.GetDialect("postgres")
	i.Searcher()
	scope := &iodata.Scope{Search: i.searcher, Dialect: dialect, Header: i.Header(), Table:i.InputName}
	s, err := scope.PrepareQuerySQL()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (ds *DataSource) query() (sql string, args []interface{}, err error) {
	scope, err := ds.NewSearchScope()
	if err != nil {
		return "", nil, err
	}
	q := &query.Query{Context: ds.Ctx, Source: scope.SQL}
	if err = q.Parse(); err != nil {
		return
	}
	c := q.NewCompiler()
	var out bytes.Buffer
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		} else {
			args = scope.SQLVars
			sql = out.String()
		}
	}()
	c.Compile(&out)
	return
}

func (ds *DataSource) Load() (api.DataReadCloser, error) {
	query, args, err := ds.query()
	if err != nil {
		return nil, err
	}
	dr := ds.ReaderFactory.Factory(ds.DataHeader, ds.db, query, args)
	return dr.(api.DataReadCloser), nil
}
