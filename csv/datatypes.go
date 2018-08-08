package csv

import (
	"strings"

	"github.com/moisespsena-go/iodata/api/datatypes"
)

func Float64(t *datatypes.Float64) *datatypes.Float64 {
	t.AddPrepareScan(func(data []byte) []byte {
		return []byte(strings.Replace(string(data), ",", ".", 1))
	})
	return t
}
