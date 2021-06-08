package nodes

type NodeStats struct {
	Name       string            `json:"name"`
	Ready      bool              `json:"ready"`
	AgeSeconds int               `json:"seconds"`
	AgeHuman   string            `json:"age"`
	LabelMap   map[string]string `json:"labels"`
}
