package datatypes

import (
	"time"

	"github.com/moisespsena-go/iodata/api"
)

func NewString() *String {
	return &String{api.NewDataTypeBase(STRING, "")}
}
func NewInt64() *Int64 {
	return &Int64{api.NewDataTypeBase(INT64, 0)}
}
func NewFloat64() *Float64 {
	return &Float64{DataTypeBase: api.NewDataTypeBase(FLOAT64, 0.0)}
}
func NewDate() *Date {
	var t time.Time
	return &Date{api.NewDataTypeBase(DATE, t)}
}

func NewTimestamp() *Date {
	var t time.Time
	return &Date{api.NewDataTypeBase(DATETIME, t)}
}

func NewBool() *Bool {
	return &Bool{api.NewDataTypeBase(BOOL, false)}
}
