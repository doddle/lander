package nodes

type NodeStats struct {
	Name        string            `json:"name"`
	Ready       bool              `json:"ready"`
	AgeSeconds  int               `json:"age"`
	Schedulable bool              `json:"schedulable"`
	LabelMap    map[string]string `json:"labels"`
	Version     string            `json:"version"`
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
