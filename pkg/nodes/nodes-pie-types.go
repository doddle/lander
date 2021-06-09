package nodes

type NodePieChart struct {
	Series []int64 `json:"series"`
}

type ChartOpts struct {
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

type FinalPieChart struct {
	ChartOpts ChartOpts `json:"chartOptions"`
	Series    []int64   `json:"series"`
	Total     int64     `json:"total"`
}

type Chart struct {
	ID         string     `json:"id"`
	DropShadow DropShadow `json:"dropShadow"`
}
type DropShadow struct {
	Effect bool `json:"effect"`
}
