package api

type DataProcessor interface {
}

type DataProcessorPlugin interface {
	Plugin
	SetInputs(map[string]DataSource)
	InputHeaders() map[string]*DataHeader
	OutputHeader() *DataHeader
	SetParameters(p map[string]interface{})
	Process(w DataWriter) error
}
