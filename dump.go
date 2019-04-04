package iodata

import (
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/error-wrap"
)

// Dump non ptr values
func Dump(header api.DataHeader, data [][]interface{}) (bytes [][][]byte, err error) {
	var (
		byts  = make([][][]byte, len(data))
		types = header.Types()
		names = header.Names()
	)

	for i, d := range data {
		byts[i] = make([][]byte, len(types))
		for j, t := range types {
			byts[i][j], err = t.Dump(d[j])
			if err != nil {
				return nil, errwrap.Wrap(err, "Dump Data[%d][%d=%v]", i, j, names[j])
			}
		}
	}
	return byts, nil
}

// DumpP dump ptr values
func DumpP(header api.DataHeader, data [][]interface{}) (bytes [][][]byte, err error) {
	var (
		byts  = make([][][]byte, len(data))
		types = header.Types()
		names = header.Names()
	)

	for i, d := range data {
		byts[i] = make([][]byte, len(types))
		for j, t := range types {
			byts[i][j], err = t.DumpP(d[j])
			if err != nil {
				return nil, errwrap.Wrap(err, "DumpP Data[%d][%d=%v]", i, j, names[j])
			}
		}
	}
	return byts, nil
}
