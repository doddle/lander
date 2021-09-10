package statefulsets

import (
	"github.com/doddle/lander/pkg/chart"
)

type FinalPieChart struct {
	ChartOpts chart.Opts `json:"chartOptions"`
	Series    []int64    `json:"series"`
	Total     int64      `json:"total"`
}
