package iodata

import (
	"io"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/error-wrap"
)

type SliceStringsReader struct {
	Slices [][]string
	Index  int
}

func (s *SliceStringsReader) Read(strings []string) error {
	if s.Index == len(s.Slices) {
		return io.EOF
	}
	for i := range strings {
		strings[i] = s.Slices[s.Index][i]
	}
	s.Index++
	return nil
}

func NewSliceStringsReader(slices ...[]string) *SliceStringsReader {
	return &SliceStringsReader{Slices: slices}
}

type StringsBytesReader struct {
	Reader api.StringsReader
	Line   int
}

func (r *StringsBytesReader) Read(result [][][]byte) (count int, err error) {
	record := make([]string, len(result[0]))
	for i := range result {
		err = r.Reader.Read(record)
		if err != nil {
			if err == io.EOF {
				return
			}
			err = errwrap.Wrap(err, "Read")
			return
		}
		for j := range result[i] {
			result[count][j] = []byte(record[j])
		}
		count++
	}
	return
}
