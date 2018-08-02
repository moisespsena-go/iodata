package iodata

type StringsBytesWriter struct {
	WriteFunc func([][]string) error
}

func (w *StringsBytesWriter) Write(data [][][]byte) (err error) {
	strings := make([][]string, len(data))
	l := len(data[0])
	for i, byts := range data {
		strings[i] = make([]string, l)
		for j, d := range byts {
			strings[i][j] = string(d)
		}
	}
	return w.WriteFunc(strings)
}
