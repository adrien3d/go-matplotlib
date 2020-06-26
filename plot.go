package main

import (
	"github.com/adrien3d/go-plotlib/utils"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
)

/*func DrawTimeChart(timestamps []time.Time, datas [][]float64, chartName string, seriesNames []string, colors []drawing.Color) {
	var timeseries []chart.Series

	for i := 0; i < len(datas); i++ {
		timeseries = append(timeseries, chart.TimeSeries{
			Name:    seriesNames[i],
			XValues: timestamps,
			YValues: datas[i],
			Style: chart.Style{
				Show:        true,
				StrokeColor: colors[i],
				StrokeWidth: 2.0,
			},
		})
	}

	graph := chart.Chart{
		DPI:    150,
		Height: 800,
		Width:  2048,
		Background: chart.Style{
			Padding:   chart.Box{Top: 50, Left: 25, Right: 25, Bottom: 10},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		XAxis: chart.XAxis{Name: "Time", NameStyle: chart.StyleShow(), Style: chart.StyleShow(), ValueFormatter: chart.TimeValueFormatterWithFormat(("00/01/02"))},
		YAxis: chart.YAxis{Name: chartName, AxisType: chart.YAxisSecondary, NameStyle: chart.StyleShow(), Style: chart.StyleShow()},
		//YAxisSecondary: chart.YAxis{Name: chartName", NameStyle: chart.StyleShow(), Style: chart.StyleShow()},
		Series: timeseries,
	}
	graph.Elements = []chart.Renderable{chart.Legend(&graph)}
	file, _ := os.Create("chart-" + chartName + ".png")
	graph.Render(chart.PNG, file)
}*/

func DrawChart(datas [][]float64, chartName string, seriesNames []string, colors []drawing.Color) {
	var series []chart.Series
	for i := 1; i < len(datas); i++ {
		series = append(series, chart.ContinuousSeries{
			Name:    seriesNames[i-1],
			XValues: datas[0],
			YValues: datas[i],
			Style: chart.Style{
				Show:        true,
				StrokeColor: colors[i],
				StrokeWidth: 2.0,
			},
		})
	}

	graph := chart.Chart{
		DPI:    150,
		Height: 600,
		Width:  1000,
		Background: chart.Style{
			Padding:   chart.Box{Top: 20, Left: 25, Right: 10, Bottom: 10},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		XAxis: chart.XAxis{Name: "Time", NameStyle: chart.StyleShow(), Style: chart.StyleShow()},
		YAxis: chart.YAxis{Name: "Values", AxisType: chart.YAxisSecondary, NameStyle: chart.StyleShow(), Style: chart.StyleShow()},
		//YAxisSecondary: chart.YAxis{Name: chartName", NameStyle: chart.StyleShow(), Style: chart.StyleShow()},
		Series: series,
	}
	graph.Elements = []chart.Renderable{chart.Legend(&graph)}
	file, _ := os.Create("chart-" + chartName + ".png")
	graph.Render(chart.PNG, file)
}

func DrawExampleChart(name string) {
	data := utils.OpenCSV(name)
	//Year,Lap time,Lap Speed
	//timestamps := make([]time.Time, 0)
	datas := make([][]float64, 3)

	for _, deviceData := range data {
		//timestamps = append(timestamps, time.Date(int(deviceData[0]), 1, 1, 0, 0, 0, 0, time.UTC))
		datas[0] = append(datas[0], deviceData[0])
		datas[1] = append(datas[1], deviceData[1])
		datas[2] = append(datas[2], deviceData[2])
	}

	lineNames := []string{"Lap time", "Speed"}
	var colors []drawing.Color
	colors = append(colors, drawing.Color{255, 0, 0, 255})     //Red
	colors = append(colors, drawing.Color{0, 255, 0, 255})     //Green
	colors = append(colors, drawing.Color{0, 0, 255, 255})     //Blue
	colors = append(colors, drawing.Color{255, 255, 255, 255}) //Black
	DrawChart(datas, name, lineNames, colors)
}
