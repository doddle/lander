package chart

/*
// These allow constructing json to tweak chart options on the front'ends piechart lib (apexcharts)
// See: https://apexcharts.com/docs/options/
Example usage:

import (
	"github.com/doddle/lander/pkg/chart"
)


type FinalPieChart struct {
	ChartOpts chart.ChartOpts `json:"chartOptions"`
	Series    []int64   `json:"series"`
	Total     int64     `json:"total"`
}

func foo(){
	result := FinalPieChart{
		Total: someNumber,
		Series: resultSeries,
		ChartOpts: chart.ChartOpts{
			Legend: chart.Legend{Show: true},
			//Theme:  chart.Theme{Palette: "palette1"},
			//Title: chart.Title{Text: "Nodes"},
			PlotOpt: chart.PlotOpt{
				Pie: chart.PlotOptPie{
					ExpandOnClick: false,
					Size:          119,
				},
			},
			Colors: resultColors,
			Stroke: chart.Stroke{Width: -1},
			Chart: chart.Chart{
				ID: "pie-nodes",
				//DropShadow: chart.DropShadow{
				//	Effect: false,
				//},
			},
			Labels: resultLabels,
		},
	}
}


*/

type Opts struct {
	Chart   Chart    `json:"chart"`
	Labels  []string `json:"labels"`
	PlotOpt PlotOpt  `json:"plotOptions"`
	Legend  Legend   `json:"legend"`
	Stroke  Stroke   `json:"stroke"`
	Theme   Theme    `json:"theme"`
	Title   Title    `json:"title"`
	Colors  []string `json:"colors"`
}

type Title struct {
	Text string `json:"text"`
}

type Data struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	BorderWidth     int64    `json:"borderWidth"`
	BorderColor     []string `json:"borderColor"`
	BackgroundColor []string `json:"backgroundColor"`
	Data            []int64  `json:"data"`
}

type Options struct {
	Responsive          bool `json:"responsive"`
	MaintainAspectRatio bool `json:"maintainAspectRatio"`
}

type PlotOpt struct {
	Pie PlotOptPie `json:"pie"`
}
type PlotOptPie struct {
	CustomScale   float64 `json:"customScale,omitempty"`
	ExpandOnClick bool    `json:"expandOnClick"`
	Size          int64   `json:"size,omitempty"`
}

type Theme struct {
	Palette string `json:"palette"`
}

type Stroke struct {
	Width float64 `json:"width"`
}

type Legend struct {
	Show bool `json:"show"`
}

type Chart struct {
	ID         string     `json:"id"`
	DropShadow DropShadow `json:"dropShadow"`
}
type DropShadow struct {
	Effect bool `json:"effect"`
}
