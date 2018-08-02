package api

type DataValue struct {
	Value    interface{}
	NotBlank bool
}

func (dv *DataValue) Blank() bool {
	return !dv.NotBlank
}
