package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/fyne-io/examples/img/icon"
	"strconv"
)

var (
	pl Plotlib
)

// Plotlib is an app to get Plotlib graph and display it
type Plotlib struct {
	image *canvas.Image
	coord *widget.Label
}

func (pl *Plotlib) downloadImage(fileName string) {
	pl.image.File = "/Users/adrien/Dev/go/src/github.com/adrien3d/go-plotlib/chart-"+fileName+".png"
	canvas.Refresh(pl.image)
}

type CustomImage struct {
	widget.Box
}

func (img *CustomImage) MouseIn(*desktop.MouseEvent) {
	fmt.Println("Entered")
}

func (img *CustomImage) MouseOut() {
	fmt.Println("Exited")
}

func (img *CustomImage) MouseMoved(me *desktop.MouseEvent) {
	//fmt.Println(me, pl.image.Position())
	pl.coord.SetText("Coord: X:" + strconv.Itoa(me.Position.X) + "    Y:" + strconv.Itoa(600 - me.Position.Y))
}

// Show starts a new Plotlib widget
func Show(app fyne.App) {
	w := app.NewWindow("Plotlib Viewer")
	w.SetIcon(icon.XKCDBitmap)

	go pl.downloadImage("indy-500-laps")
	submit := widget.NewButton("Submit", func() {
		//x.Submit()
	})
	submit.Style = widget.PrimaryButton
	pl.coord = widget.NewLabel("Coordinates")
	pl.coord.Alignment = fyne.TextAlignCenter
	buttonsBox := widget.NewHBox(
		widget.NewButton("Exit", func() {
			w.Close()
		}),
		layout.NewSpacer(),
		pl.coord,
		submit)
	pl.image = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	imageBox := widget.NewHBox(pl.image)
	content := &CustomImage{*imageBox}
	vbox := widget.NewVBox(content, buttonsBox)
	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), vbox))
	w.Show()
}

func main() {
	DrawExampleChart("indy-500-laps")
	launch := Show
	ex := app.New()
	launch(ex)
	ex.Run()
}
