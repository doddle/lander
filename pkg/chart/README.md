
# Chart package

These types allow constructing json to tweak chart options on the front'ends piechart lib


## Example usage:


```

package foo

import (
    "github.com/digtux/lander/pkg/chart"
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


```
