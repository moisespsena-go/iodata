package iodata

import (
	"fmt"
	"testing"
)

func TestParseParameters(t *testing.T) {
	params := struct {
		A int
		B struct {
			C float64
			D struct {
				E bool
			}
		}
		X float64
		Y []string
	}{}
	a, c, e, x, y := 15, 4.5, true, 7.9, []string{"1", "2"}
	ParseParameters(map[string]interface{}{
		"A":     a,
		"B.C":   c,
		"B.D.E": e,
		"X":     x,
		"Y":     y,
	}, &params)

	if params.A != a {
		t.Fatalf("Expected <%v>, but get %q", a, params.A)
	}
	if params.B.C != c {
		t.Fatalf("Expected <%v>, but get %q", c, params.B.C)
	}
	if params.B.D.E != e {
		t.Fatalf("Expected <%v>, but get %q", e, params.B.D.E)
	}
	if params.X != x {
		t.Fatalf("Expected <%v>, but get %q", x, params.X)
	}
	if fmt.Sprint(params.Y) != fmt.Sprint(y) {
		t.Fatalf("Expected <%v>, but get %q", y, params.Y)
	}
}
