package modelstruct

import "testing"

type a struct {
	f1 int
	f2 float64
}

func TestGet(t *testing.T) {
	m := Get(&a{})
	if m == nil {
		t.Fatal("model is nil")
	}
	if len(m.Fields) != 2 {
		t.Fatal("num fields != 2")
	}

	checkField := func(i int, name string) {
		if m.Fields[i].Name != name {
			t.Fatalf("Field %d: name != %q", i, name)
		}
	}

	for i, name := range []string{"f1", "f2"} {
		checkField(i, name)
	}
}

type b struct {
	a
	x int
}

func TestGet_Anonymous(t *testing.T) {
	m := Get(&b{})
	if m == nil {
		t.Fatal("model is nil")
	}
	if len(m.Fields) != 3 {
		t.Fatal("num fields != 3")
	}

	checkField := func(i int, name string) {
		if m.Fields[i].Name != name {
			t.Fatalf("Field %d: name != %q", i, name)
		}
	}

	for i, name := range []string{"f1", "f2", "x"} {
		checkField(i, name)
	}
}