package db

import (
	"database/sql"

	"github.com/moisespsena-go/iodata/api"
)

func NewReaderFactory() api.DataReaderFactory {
	return api.DataReaderFactoryFunc(func(header api.DataHeader, args ...interface{}) api.DataReader {
		var (
			db      = args[0].(*sql.DB)
			SQL     = args[1].(string)
			sqlArgs = args[2].([]interface{})
		)
		return &Reader{DB: db, SQL: SQL, SQLArgs: sqlArgs}
	})
}

func NewWriterFactory() api.DataWriterFactory {
	return api.DataWriterFactoryFunc(func(header api.DataHeader, args ...interface{}) api.DataWriter {
		return nil
	})
}
