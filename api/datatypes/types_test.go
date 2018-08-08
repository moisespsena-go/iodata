package datatypes

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"
)

func TestString_Scan(t *testing.T) {
	expectedValue := "value"
	dt := String{}
	v, err := dt.Scan([]byte(expectedValue))
	if err != nil {
		t.Error(err)
	}
	if v.(string) != expectedValue {
		t.Fail()
	}
}

func TestString_Dump(t *testing.T) {
	expectedValue := "value"
	dt := String{}
	v, err := dt.Dump(expectedValue)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fail()
	}
}

func TestInt64_Scan(t *testing.T) {
	var expectedValue int64 = 1998
	dt := Int64{}
	v, err := dt.Scan([]byte("1998"))
	if err != nil {
		t.Error(err)
	}
	if v.(int64) != expectedValue {
		t.Fail()
	}
}

func TestInt64_Dump(t *testing.T) {
	expectedValue := "1998"
	dt := Int64{}
	v, err := dt.Dump(1998)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fail()
	}
}

func TestFloat64_Scan(t *testing.T) {
	expectedValue := 1998.78
	dt := Float64{}
	v, err := dt.Scan([]byte("1998.78"))
	if err != nil {
		t.Error(err)
	}
	if v.(float64) != expectedValue {
		t.Fail()
	}
	expectedValue = 1e-9
	v, err = dt.Scan([]byte("1e-9"))
	if err != nil {
		t.Error(err)
	}
	if v.(float64) != expectedValue {
		t.Fail()
	}
	v, err = dt.Scan([]byte("1E-9"))
	if err != nil {
		t.Error(err)
	}
	if v.(float64) != expectedValue {
		t.Fail()
	}
}

func TestFloat64_ScanRound(t *testing.T) {
	expectedValue := 1998.123
	dt := Float64{Round: &FloatRound{Places: 3}}
	v, err := dt.Scan([]byte("1998.123456"))
	if err != nil {
		t.Error(err)
	}
	if v.(float64) != float64(expectedValue) {
		t.Fatalf("Expected <%v>, but again <%v>", expectedValue, v)
	}

	expectedValue = 1998.124
	dt.Round.RoundFunc = math.Ceil
	v, err = dt.Scan([]byte("1998.123456"))
	if err != nil {
		t.Error(err)
	}
	if v.(float64) != float64(expectedValue) {
		t.Fatalf("Expected <%v>, but get <%v>", expectedValue, v)
	}
}

func TestFloat64_Dump(t *testing.T) {
	expectedValue := "1998.78"
	dt := Float64{}
	v, err := dt.Dump(1998.78)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fail()
	}

	dt.Round = &FloatRound{Places: 2}
	v, err = dt.Dump(-1.234456E+78)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte("-1.234456e+78")) != 0 {
		t.Fatalf("Expected <%v>, but get <%v>", "-1.234456e+78", string(v))
	}

	dt.Format = "%E"
	v, err = dt.Dump(-1.234456E+78)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte("-1.234456E+78")) != 0 {
		t.Fatalf("Expected <%v>, but get <%v>", "-1.234456E+78", string(v))
	}
}

func TestBool_Scan(t *testing.T) {
	dt := Bool{}

	for _, v := range []string{"true", "1", "t", "y"} {
		v, err := dt.Scan([]byte(v))
		if err != nil {
			t.Error(err)
		}
		if !v.(bool) {
			t.Errorf("Scan %q as true failed.", v)
		}
	}

	for _, v := range []string{"f", "false", "0"} {
		v, err := dt.Scan([]byte(v))
		if err != nil {
			t.Error(err)
		}
		if v.(bool) {
			t.Errorf("Scan %q as false failed.", v)
		}
	}

	_, err := dt.Scan([]byte("-"))
	if err == nil {
		t.Error("Error expected")
	}
}

func TestBool_Dump(t *testing.T) {
	expectedValue := "true"
	dt := Bool{}
	v, err := dt.Dump(true)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fail()
	}
}

func TestDate_Scan(t *testing.T) {
	dt := Date{}
	expected := time.Date(2018, time.Month(07), 30, 0, 0, 0, 0, &time.Location{})

	v, err := dt.Scan([]byte("2018-07-30"))
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprint(expected) != fmt.Sprint(v) {
		t.Fatalf("Expected <%v>, but get <%v>", expected, v)
	}
}

func TestDate_Dump(t *testing.T) {
	expectedValue := "2018-07-30"
	dt := Date{}
	v, err := dt.Dump(time.Date(2018, time.Month(07), 30, 0, 0, 0, 0, &time.Location{}))
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fatalf("Expected <%v>, but get <%v>", expectedValue, string(v))
	}
}

func TestTimestamp_Scan(t *testing.T) {
	dt := Timestamp{}
	expected := time.Date(2018, time.Month(07), 30, 13, 22, 59, 0, &time.Location{})

	v, err := dt.Scan([]byte("2018-07-30 13:22:59"))
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprint(expected) != fmt.Sprint(v) {
		t.Fatalf("Expected <%v>, but get <%v>", expected, v)
	}
}

func TestTimestamp_Dump(t *testing.T) {
	expectedValue := "2018-07-30 13:22:59"
	dt := Timestamp{}
	v, err := dt.Dump(time.Date(2018, time.Month(07), 30, 13, 22, 59, 0, &time.Location{}))
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(v, []byte(expectedValue)) != 0 {
		t.Fatalf("Expected <%v>, but get <%v>", expectedValue, string(v))
	}
}
