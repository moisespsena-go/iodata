package iodata

import (
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/api/datatypes"
)

func NewDataHeader(args ...interface{}) api.DataHeader {
	if len(args) == 1 {
		switch at := args[0].(type) {
		case []Field:
			return api.NewDataHeader(DataFields(datatypes.DefaultTypes, at...)...)
		}
	}
	return nil
}

type Field struct {
	Name string
	Type api.DataTypeName
}

func DataFields(types map[api.DataTypeName]api.DataType, field ...Field) []api.DataField {
	dataFields := make([]api.DataField, len(field))

	var (
		typeName api.DataTypeName
		tn       string
		ok       bool
		t        api.DataType
	)

	for i, f := range field {
		typeName = f.Type
		tn = string(typeName.Name())
		if t, ok = types[api.DataTypeName(tn)]; !ok {
			panic("DataType \"" + tn + "\" not exists.")
		}
		dataFields[i] = api.NewDataField(f.Name, t)
	}
	return dataFields
}
