package db

import (
	"reflect"
	"testing"
)

func TestRawScaner_Scan(t *testing.T) {
	var v int64 = 355
	vs := &RawScaner{
		Typ:  reflect.TypeOf(v),
		Data: &v,
	}

	// test with value
	if err := vs.Scan("100"); err != nil {
		t.Fatal(err)
	}
	if !vs.NotNil {
		t.Fatal("Expected NotNil.")
	}
	if *vs.Data.(*int64) != 100 {
		t.Fatal("Expected 100.")
	}
	if v != 355 {
		t.Fatal("Variable v changed.")
	}

	// test nil
	vs = &RawScaner{
		Typ:  reflect.TypeOf(v),
		Data: &v,
	}
	if err := vs.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if vs.NotNil {
		t.Fatal("Expected Nil.")
	}
	if *vs.Data.(*int64) != 355 {
		t.Fatal("Expected 355.")
	}
	if v != 355 {
		t.Fatal("Variable v changed.")
	}
}
