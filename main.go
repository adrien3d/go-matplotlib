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

	currentPoint    = -1
	lineLength      = 10
	inputData       [][]float64
	dataLimits      [5]float64 //xMin, xMax, yMin, yMax, nValues
	datasetNames    []string
	datasetSelected []int64
	//graphDimensions [2]int64
)

// Margins
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

const (
	colorWhite      = 0xFFFFFF
	colorBlack      = 0x000000
	colorDodgerBlue = 0x1E90FF
)

func getData(filename string) [][]float64 {
	columnNames, csvData := utils.OpenCSV(filename)
	datasetNames = columnNames
	datas := make([][]float64, len(csvData[0]))

	for i := 0; i < len(csvData); i++ { //all lines
		//timestamps = append(timestamps, time.Date(int(deviceData[0]), 1, 1, 0, 0, 0, 0, time.UTC))
		datas[0] = append(datas[0], csvData[i][0])
		datas[1] = append(datas[1], csvData[i][1])
		datas[2] = append(datas[2], csvData[i][2])
	}
	return datas
}

func refreshLimits() {
	dataLimits = [5]float64{inputData[0][0], inputData[0][len(inputData[0])-1], inputData[1][0], inputData[1][len(inputData[1])-1], float64(len(inputData[0]))}
	for _, dsSelectedIndex := range datasetSelected { // Up to 3
		for j := 0; j < len(inputData[dsSelectedIndex]); j++ { // All file length (i.e. 70)
			//X
			if inputData[0][j] < dataLimits[0] {
				dataLimits[0] = inputData[0][j]
			}
			if inputData[0][j] > dataLimits[1] {
				dataLimits[1] = inputData[0][j]
			}
			//Y
			if inputData[dsSelectedIndex][j] < dataLimits[2] {
				dataLimits[2] = inputData[dsSelectedIndex][j]
			}
			if inputData[dsSelectedIndex][j] > dataLimits[3] {
				dataLimits[3] = inputData[dsSelectedIndex][j]
			}
		}
	}
}

func pointLocations(areaWidth, areaHeight float64) (xs, ys [100]float64) {
	lineLength = int(dataLimits[4])
	xincr := areaWidth / (dataLimits[1] - dataLimits[0])  // Or float64(lineLength) for 0 X ?
	yincr := areaHeight / (dataLimits[3] - dataLimits[2]) // Or float64(lineLength) for 0 Y ?
	for i := 0; i < lineLength; i++ {
		xs[i] = xincr * float64(i)
		ys[i] = (yincr * inputData[1][i]) - yincr*dataLimits[2] + (yoffBottom - yoffTop)
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
		path.LineTo(xs[i], graphHeight-ys[i])
	}

	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now draw the point being hovered over
	if currentPoint != -1 {
		path = ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(xs[currentPoint], graphHeight-ys[currentPoint], pointRadius, 0, 6.23, false)
		//TODO: detect which grah we are hovering
		positionLabel.SetText("X:" + fmt.Sprintf("%f", inputData[0][currentPoint]) + "\t Y:" + fmt.Sprintf("%f", inputData[datasetSelected[0]][currentPoint]))
		path.End()
		// use the same brush as for the histogram lines
		p.Context.Fill(path, brush)
		path.Free()
	}
}

func inPoint(x, y float64, xtest, ytest float64) bool {
	// TODO switch to using a matrix
	x -= xoffLeft
	y -= yoffBottom //(yoffBottom - yoffTop)
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

func getColumnIfFromDatasetName(datasetName string) (ret int64) {
	return ret
}

func setupUI() {
	inputData = getData("indy-500-laps")
	vBoxFiles := ui.NewVerticalBox()
	for i := 1; i < len(datasetNames); i++ {
		cb := ui.NewCheckbox(datasetNames[i])
		if i == 1 {
			cb.SetChecked(true)
			datasetSelected = append(datasetSelected, 1)
		}
		cb.OnToggled(func(checkbox *ui.Checkbox) {
			//TODO: retrieve i and add/remove it from datasetSelected, then refreshLimits(), then redraw selected
			if cb.Checked() {
				fmt.Println(cb.Text(), "checked")
			} else {
				fmt.Println(cb.Text(), "unchecked")
			}
			fmt.Println(dataLimits)
		})
		vBoxFiles.Append(cb, false)
	}

	refreshLimits()
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

	vboxMain := ui.NewVerticalBox()

	histogram = ui.NewArea(areaHandler{})

	hboxGraphAndFiles := ui.NewHorizontalBox()
	hboxGraphAndFiles.Append(histogram, false)
	hboxGraphAndFiles.Append(vBoxFiles, false)

	hboxTools := ui.NewHorizontalBox()
	resetButton = ui.NewButton("Home")
	positionLabel = ui.NewLabel("X: \t Y:")
	hboxTools.Append(resetButton, false)
	hboxTools.Append(positionLabel, false)

	vboxMain.Append(hboxGraphAndFiles, false)
	vboxMain.Append(hboxTools, false)

	mainwin.SetChild(vboxMain)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
