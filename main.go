package main

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/fyne-io/examples/img/icon"
)

var (
	pl Plotlib
)

// Plotlib is an app to get Plotlib graph and display it
type Plotlib struct {
	image *canvas.Image
	coord *widget.Label
}

func (pl *Plotlib) downloadImage(url string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile(os.TempDir(), "Plotlib.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	pl.image.File = file.Name()
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
	pl.coord.SetText("Coord: X:" + strconv.Itoa(me.Position.X) + "    Y:" + strconv.Itoa(me.Position.Y))
}

// Show starts a new Plotlib widget
func Show(app fyne.App) {
	w := app.NewWindow("Plotlib Viewer")
	w.SetIcon(icon.XKCDBitmap)

	go pl.downloadImage("https://www.google.fr/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png")
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
	launch := Show
	ex := app.New()
	launch(ex)
	ex.Run()
}
