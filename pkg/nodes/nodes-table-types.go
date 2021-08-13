package nodes

// Types for table data

type NodeStats struct {
	Name       string            `json:"name"`
	Ready      bool              `json:"ready"`
	AgeSeconds int               `json:"seconds"`
	AgeHuman   string            `json:"age"`
	LabelMap   map[string]string `json:"labels"`
}
type NodeTable struct {
	Nodes   []NodeStats    `json:"nodes"`
	Headers []TableHeaders `json:"headers"`
}

type TableHeaders struct {
	Align string `json:"align,omitempty"`
	Text  string `json:"text"`
	Value string `json:"value"`
}
