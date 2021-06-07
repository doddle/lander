package pie_statefulset

import "encoding/json"

func UnmarshalWelcome(data []byte) (StatefulSetPieChart, error) {
	var r StatefulSetPieChart
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StatefulSetPieChart) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StatefulSetPieChart struct {
	Series []int64 `json:"series"`
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	BorderWidth     int64    `json:"borderWidth"`
	BorderColor     []string `json:"borderColor"`
	BackgroundColor []string `json:"backgroundColor"`
	Data            []int64  `json:"data"`
}

type ChartOptions struct {
	Legend              Legend `json:"legend"`
	Responsive          bool   `json:"responsive"`
	MaintainAspectRatio bool   `json:"maintainAspectRatio"`
}

type Legend struct {
	Display bool `json:"display"`
}
