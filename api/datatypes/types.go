package datatypes

import (
	"fmt"
	"strconv"

	"reflect"

	"math"

	"regexp"

	"time"

	"github.com/moisespsena-go/iodata/api"
)

type DataTypeName = api.DataTypeName

type DataType struct {
	api.DataTypeBase
	ScanFunc func(value []byte) (interface{}, error)
	DumpFunc func(value interface{}) ([]byte, error)
}

func (d DataType) Scan(value []byte) (interface{}, error) {
	return d.ScanFunc(value)
}

func (d DataType) Dump(value interface{}) ([]byte, error) {
	return d.DumpFunc(value)
}

const (
	STRING   DataTypeName = "string"
	INT64    DataTypeName = "int64"
	FLOAT64  DataTypeName = "float64"
	DATE     DataTypeName = "date"
	DATETIME DataTypeName = "datetime"
	BOOL     DataTypeName = "bool"
)

var (
	StringReflectType  = reflect.TypeOf((*string)(nil)).Elem()
	Int64ReflectType   = reflect.TypeOf((*int64)(nil)).Elem()
	Float64ReflectType = reflect.TypeOf((*float64)(nil)).Elem()
	BoolReflectType    = reflect.TypeOf((*bool)(nil)).Elem()
	TimeReflectType    = reflect.TypeOf((*time.Time)(nil)).Elem()
	DefaultTypes       = map[DataTypeName]api.DataType{
		STRING:   NewString(),
		INT64:    NewInt64(),
		FLOAT64:  NewFloat64(),
		DATE:     NewDate(),
		DATETIME: NewTimestamp(),
		BOOL:     NewBool(),
	}
)

type String struct {
	api.DataTypeBase
}

func (String) Scan(value []byte) (interface{}, error) {
	return string(value), nil
}

func (String) Dump(value interface{}) ([]byte, error) {
	return []byte(value.(string)), nil
}

func (t String) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*string))
}

func (String) Elem() api.DataType {
	return nil
}

func (String) Type() reflect.Type {
	return StringReflectType
}

type Int64 struct {
	api.DataTypeBase
}

func (Int64) Scan(value []byte) (interface{}, error) {
	v, err := strconv.Atoi(string(value))
	return int64(v), err
}

func (Int64) Dump(value interface{}) ([]byte, error) {
	return []byte(fmt.Sprint(value.(int64))), nil
}

func (t Int64) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*int64))
}

func (Int64) Elem() api.DataType {
	return nil
}

func (Int64) Type() reflect.Type {
	return Int64ReflectType
}

type FloatRound struct {
	Places    int
	RoundFunc func(value float64) float64
}

func (f FloatRound) round(value float64) int {
	return int(value + math.Copysign(0.5, value))
}

func (f FloatRound) DefaultFix(value float64) float64 {
	output := math.Pow(10, float64(f.Places))
	return float64(f.round(value*output)) / output
}

func (f FloatRound) Fix(value float64) float64 {
	if f.RoundFunc != nil {
		exp := math.Pow10(f.Places)
		value = f.RoundFunc(value*exp) / exp
		return value
	}
	return f.DefaultFix(value)
}

type Float64 struct {
	api.DataTypeBase
	Format string
	Round  *FloatRound
}

func (f Float64) Scan(value []byte) (interface{}, error) {
	v, err := strconv.ParseFloat(string(value), 64)
	if err == nil && f.Round != nil {
		v = f.Round.Fix(v)
	}
	return v, err
}

func (f Float64) Dump(value interface{}) ([]byte, error) {
	fvalue := value.(float64)
	if f.Format == "" {
		return []byte(fmt.Sprint(fvalue)), nil
	}
	return []byte(fmt.Sprintf(f.Format, fvalue)), nil
}

func (t Float64) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*float64))
}

func (Float64) Elem() api.DataType {
	return nil
}

func (Float64) Type() reflect.Type {
	return Float64ReflectType
}

type Bool struct {
	api.DataTypeBase
}

func (Bool) Scan(value []byte) (interface{}, error) {
	var (
		v   bool
		err error
	)

	vs := string(value)
	switch vs {
	case "t", "true", "1", "y":
		v = true
		break
	case "f", "false", "0":
		break
	default:
		err = fmt.Errorf("invalid bool value: %q", vs)
	}
	return v, err
}

func (Bool) Dump(value interface{}) ([]byte, error) {
	v := "false"
	if value.(bool) {
		v = "true"
	}
	return []byte(v), nil
}

func (t Bool) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*bool))
}

func (Bool) Elem() api.DataType {
	return nil
}

func (Bool) Type() reflect.Type {
	return BoolReflectType
}

type Date struct {
	api.DataTypeBase
}

var (
	dateRe, _ = regexp.Compile("(\\d{4})-(\\d{2})-(\\d{2})")
)

func (Date) Scan(value []byte) (interface{}, error) {
	var (
		v time.Time
	)

	r := dateRe.FindAllStringSubmatch(string(value), 1)
	if len(r) == 0 {
		return v, fmt.Errorf("invalid date value: %q", string(value))
	}
	y, _ := strconv.Atoi(r[0][1])
	m, _ := strconv.Atoi(r[0][2])
	d, _ := strconv.Atoi(r[0][3])
	v = time.Date(y, time.Month(m), d, 0, 0, 0, 0, &time.Location{})
	return v, nil
}

func (Date) Dump(value interface{}) ([]byte, error) {
	t := value.(time.Time)
	return []byte(fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())), nil
}

func (t Date) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*time.Time))
}

func (Date) Elem() api.DataType {
	return nil
}

func (Date) Type() reflect.Type {
	return TimeReflectType
}

type Timestamp struct {
	api.DataTypeBase
}

var (
	timestampRe, _ = regexp.Compile("(\\d{4})-(\\d{2})-(\\d{2}) (\\d{2}):(\\d{2}):(\\d{2})")
)

func (Timestamp) Scan(value []byte) (interface{}, error) {
	var (
		v time.Time
	)

	r := timestampRe.FindAllStringSubmatch(string(value), 1)
	if len(r) == 0 {
		return v, fmt.Errorf("invalid timestamp value: %q", string(value))
	}
	y, _ := strconv.Atoi(r[0][1])
	m, _ := strconv.Atoi(r[0][2])
	d, _ := strconv.Atoi(r[0][3])
	h, _ := strconv.Atoi(r[0][4])
	min, _ := strconv.Atoi(r[0][5])
	sec, _ := strconv.Atoi(r[0][6])
	v = time.Date(y, time.Month(m), d, h, min, sec, 0, &time.Location{})
	return v, nil
}

func (Timestamp) Dump(value interface{}) ([]byte, error) {
	t := value.(time.Time)
	return []byte(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())), nil
}

func (t Timestamp) DumpP(value interface{}) ([]byte, error) {
	return t.Dump(*value.(*time.Time))
}

func (Timestamp) Elem() api.DataType {
	return nil
}

func (Timestamp) Type() reflect.Type {
	return TimeReflectType
}
