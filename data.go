package iodata

import "github.com/moisespsena-go/iodata/api"

type Field struct {
	Name string
	Type string
}

func DataFields(types map[string]api.DataType, field ...Field) []api.DataField {
	dataFields := make([]api.DataField, len(field))

	var (
		typeName api.DataTypeName
		tn       string
		ok       bool
		t        api.DataType
	)

	for i, f := range field {
		typeName = api.DataTypeName(f.Type)
		tn = string(typeName.Name())
		if t, ok = types[tn]; !ok {
			panic("DataType \"" + tn + "\" not exists.")
		}
		dataFields[i] = api.NewDataField(f.Name, t)
	}
	return dataFields
}
