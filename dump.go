package iodata

import (
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

func Dump(header *api.DataHeader, data [][]interface{}) (bytes [][][]byte, err error) {
	byts := make([][][]byte, len(data))

	for i, d := range data {
		for j, t := range header.Types {
			byts[i] = make([][]byte, len(header.Types))
			byts[i][j], err = t.Dump(d[j])
			if err != nil {
				return nil, errwrap.Wrap(err, "Dump Data[%d][%d=%v]", i, j, header.Names[j])
			}
		}
	}
	return byts, nil
}
