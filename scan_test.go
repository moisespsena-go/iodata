package iodata

import (
	"testing"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/iodata/api/datatypes"
)

func TestScanAssign(t *testing.T) {
	var (
		dt       = &datatypes.Int64Type{}
		value    int64
		notBlank bool
	)

	notBlank, err := ScanAssign(dt, []byte("2018"), &value)
	if err != nil {
		t.Error(err)
	}
	if value != 2018 {
		t.Fatalf("Expected <%v>, but get %q", 2018, value)
	}

	value = 0
	if notBlank, err = ScanAssign(dt, []byte(""), &value); err != nil {
		t.Error(err)
	}
	if notBlank {
		t.Fatalf("Expected BLANK")
	}

	value = 0
	dt.Default = int64(100)
	if notBlank, err = ScanAssign(dt, []byte(""), &value); err != nil {
		t.Error(err)
	}
	if notBlank {
		t.Fatalf("Expected BLANK")
	}
	if value != dt.Default {
		t.Fatalf("Expected <%v>, but get %q", dt.Default, value)
	}
}

func TestScan(t *testing.T) {
	dt := &datatypes.Int64Type{}
	value, err := Scan(dt, []byte("2018"))
	if err != nil {
		t.Error(err)
	}
	if value.Value.(int64) != 2018 {
		t.Fatalf("Expected <%v>, but get %q", 2018, value.Value)
	}

	if value, err = Scan(dt, []byte("")); err != nil {
		t.Error(err)
	}
	if value.NotBlank {
		t.Fatalf("Expected BLANK")
	}

	dt.Default = int64(100)
	if value, err = Scan(dt, []byte("")); err != nil {
		t.Error(err)
	}
	if value.NotBlank {
		t.Fatalf("Expected BLANK")
	}
	if value.Value.(int64) != dt.Default {
		t.Fatalf("Expected <%v>, but get %q", dt.Default, value.Value)
	}
}

func TestScanSliceAssign(t *testing.T) {
	var (
		intValue   int64
		floatValue float64
	)
	_, err := ScanSliceAssign([]api.DataType{&datatypes.Int64Type{}, &datatypes.Float64Type{}},
		[][]byte{[]byte("2018"), []byte("12.8")}, []interface{}{&intValue, &floatValue})
	if err != nil {
		t.Error(err)
	}
	if intValue != 2018 {
		t.Fatalf("Expected <%v>, but get %q", 2018, intValue)
	}
	if floatValue != 12.8 {
		t.Fatalf("Expected <%v>, but get %q", 12.8, intValue)
	}
}

func TestScanSliceStringsAssign(t *testing.T) {
	var (
		intValue   int64
		floatValue float64
	)
	_, err := ScanSliceStringsAssign([]api.DataType{&datatypes.Int64Type{}, &datatypes.Float64Type{}},
		[]string{"2018", "12.8"}, []interface{}{&intValue, &floatValue})
	if err != nil {
		t.Error(err)
	}
	if intValue != 2018 {
		t.Fatalf("Expected <%v>, but get %q", 2018, intValue)
	}
	if floatValue != 12.8 {
		t.Fatalf("Expected <%v>, but get %q", 12.8, intValue)
	}
}
