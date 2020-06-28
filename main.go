package main

import (
	"fmt"
	"github.com/adrien3d/go-plotlib/utils"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	histogram     *ui.Area
	resetButton   *ui.Button
	positionLabel *ui.Label

	currentPoint = -1
	lineLength = 10
)

// some metrics
const (
	xoffLeft    = 30 // histogram margins
	yoffTop     = 20
	xoffRight   = 20
	yoffBottom  = 30
	pointRadius = 5
)

// helper to quickly set a brush color
func mkSolidBrush(color uint32, alpha float64) *ui.DrawBrush {
	brush := new(ui.DrawBrush)
	brush.Type = ui.DrawBrushTypeSolid
	component := uint8((color >> 16) & 0xFF)
	brush.R = float64(component) / 255
	component = uint8((color >> 8) & 0xFF)
	brush.G = float64(component) / 255
	component = uint8(color & 0xFF)
	brush.B = float64(component) / 255
	brush.A = alpha
	return brush
}

// and some colors
const (
	colorWhite      = 0xFFFFFF
	colorBlack      = 0x000000
	colorDodgerBlue = 0x1E90FF
)

func getData(filename string) [][]float64 {
	csvData := utils.OpenCSV(filename)
	datas := make([][]float64, len(csvData[0]))

	for _, deviceData := range csvData {
		//timestamps = append(timestamps, time.Date(int(deviceData[0]), 1, 1, 0, 0, 0, 0, time.UTC))
		datas[0] = append(datas[0], deviceData[0])
		datas[1] = append(datas[1], deviceData[1])
		datas[2] = append(datas[2], deviceData[2])
	}
	return datas
}

func pointLocations(width, height float64) (xs, ys [100]float64) {
	data := getData("indy-500-laps")
	lineLength = len(data[0])
	xincr := width / float64(lineLength) - 1 // 10 - 1 to make the last point be at the end
	//yincr := height / 10
	for i := 0; i < lineLength; i++ {
		// get the value of the point
		n := data[1][i]
		// because y=0 is the top but n=0 is the bottom, we need to flip
		n = 100 - n
		xs[i] = xincr * float64(i)
		ys[i] = float64(n)
	}
	return xs, ys
}

func graphSize(clientWidth, clientHeight float64) (graphWidth, graphHeight float64) {
	return clientWidth - xoffLeft - xoffRight, clientHeight - yoffTop - yoffBottom
}

type areaHandler struct{}

func (areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	// fill the area with white
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, p.AreaWidth, p.AreaHeight)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()

	graphWidth, graphHeight := graphSize(p.AreaWidth, p.AreaHeight)

	sp := &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}

	// draw the axes
	brush = mkSolidBrush(colorBlack, 1.0)
	path = ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(xoffLeft, yoffTop)
	path.LineTo(xoffLeft, yoffTop+graphHeight)
	path.LineTo(xoffLeft+graphWidth, yoffTop+graphHeight)
	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now transform the coordinate space so (0, 0) is the top-left corner of the graph
	m := ui.DrawNewMatrix()
	m.Translate(xoffLeft, yoffTop)
	p.Context.Transform(m)

	brush = mkSolidBrush(colorDodgerBlue, 1.0)
	// now draw the graph line
	xs, ys := pointLocations(graphWidth, graphHeight)
	path = ui.DrawNewPath(ui.DrawFillModeWinding)

	path.NewFigure(xs[0], ys[0])
	for i := 1; i < lineLength; i++ {
		path.LineTo(xs[i], ys[i])
	}

	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now draw the point being hovered over
	if currentPoint != -1 {
		xs, ys := pointLocations(graphWidth, graphHeight)
		path = ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(
			xs[currentPoint], ys[currentPoint],
			pointRadius,
			0, 6.23, // TODO pi
			false)
		path.End()
		// use the same brush as for the histogram lines
		p.Context.Fill(path, brush)
		path.Free()
	}
}

func inPoint(x, y float64, xtest, ytest float64) bool {
	// TODO switch to using a matrix
	x -= xoffLeft
	y -= yoffTop
	return (x >= xtest-pointRadius) &&
		(x <= xtest+pointRadius) &&
		(y >= ytest-pointRadius) &&
		(y <= ytest+pointRadius)
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	graphWidth, graphHeight := graphSize(me.AreaWidth, me.AreaHeight)
	xs, ys := pointLocations(graphWidth, graphHeight)

	currentPoint = -1
	for i := 0; i < lineLength; i++ {
		if inPoint(me.X, me.Y, xs[i], ys[i]) {
			currentPoint = i
			break
		}
	}

	positionLabel.SetText("X:" + fmt.Sprintf("%f", me.X) + "\t Y:" + fmt.Sprintf("%f", me.AreaHeight-me.Y))

	// TODO only redraw the relevant area
	histogram.QueueRedrawAll()
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (areaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func setupUI() {
	mainwin := ui.NewWindow("Figure 1", 640, 480, true)
	mainwin.SetMargined(false)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	mainwin.SetChild(vbox)

	histogram = ui.NewArea(areaHandler{})
	resetButton = ui.NewButton("Home")

	vbox.Append(histogram, false)
	vbox.Append(ui.NewVerticalSeparator(), false)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(resetButton, false)
	positionLabel = ui.NewLabel("X: \t Y:")
	hbox.Append(positionLabel, false)
	vbox.Append(hbox, false)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
