package statefulsets

import "encoding/json"

func (r *StatefulSetPieChart) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StatefulSetPieChart struct {
	Series []int64 `json:"series"`
}

type ChartOpts struct {
	Chart Chart 	`json:"chart"`
	Colors []string `json:"colors"`
	Labels []string `json:"labels"`
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

type FinalResult struct {
	ChartOpts ChartOpts `json:"chartOptions"`
	Series 	  []int64   `json:"series"`
	Total 	  int64     `json:"total"`
}

type Chart struct {
	ID string `json:"id"`
}

