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
	lineLength   = 10
	inputData    [][]float64
	dataLimits   [5]float64 //xMin, xMax, yMin, yMax, nValues
	//graphDimensions [2]int64
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

	dataLimits = [5]float64{csvData[0][0], csvData[len(csvData)-1][0], csvData[0][1], csvData[0][1], float64(len(csvData))}
	for i := 0; i < len(csvData); i++ {
		//timestamps = append(timestamps, time.Date(int(deviceData[0]), 1, 1, 0, 0, 0, 0, time.UTC))
		datas[0] = append(datas[0], csvData[i][0])
		datas[1] = append(datas[1], csvData[i][1])
		datas[2] = append(datas[2], csvData[i][2])
		if csvData[i][1] > dataLimits[3] {
			dataLimits[3] = csvData[i][1]
		} else if csvData[i][1] < dataLimits[2] {
			dataLimits[2] = csvData[i][1]
		}
	}
	return datas
}

func pointLocations(areaWidth, areaHeight float64) (xs, ys [100]float64) {
	lineLength = int(dataLimits[4])
	xincr := areaWidth / float64(lineLength) // - 1 to make the last point be at the end
	yincr := areaHeight / (dataLimits[3] - dataLimits[2])
	for i := 0; i < lineLength; i++ {
		xs[i] = xincr * float64(i)
		ys[i] = (yincr * inputData[1][i]) - yincr*dataLimits[2]
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
	m.Translate(xoffLeft, yoffBottom)
	p.Context.Transform(m)

	brush = mkSolidBrush(colorDodgerBlue, 1.0)
	// now draw the graph line
	xs, ys := pointLocations(graphWidth, graphHeight)
	path = ui.DrawNewPath(ui.DrawFillModeWinding)

	path.NewFigure(xs[0], graphHeight-ys[0])
	for i := 1; i < lineLength; i++ {
		path.LineTo(xs[i], graphHeight-10-ys[i])
	}

	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now draw the point being hovered over
	if currentPoint != -1 {
		path = ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(xs[currentPoint], graphHeight-ys[currentPoint]-10, pointRadius, 0, 6.23, false)
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
		if inPoint(me.X, me.Y, xs[i], graphHeight-ys[i]) {
			currentPoint = i
			break
		}
	}
	properX, properY := dataLimits[0]+(me.X-xoffLeft)*((dataLimits[1]-dataLimits[0])/graphWidth), dataLimits[2]+(graphHeight-me.Y+yoffTop)*((dataLimits[3]-dataLimits[2])/graphHeight)

	//TODO: Scale from relative position to graph values
	if properX < dataLimits[0] {
		properX = dataLimits[0]
	}
	if properY < dataLimits[2] {
		properY = dataLimits[2]
	}

	positionLabel.SetText("X:" + fmt.Sprintf("%f", properX) + "\t Y:" + fmt.Sprintf("%f", properY))

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
	inputData = getData("indy-500-laps")
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
