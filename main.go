package main

import (

	// STDLIB
	"fmt"
	"os"

	// External
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	// Internal
	"github.com/marshall-mcmullen/wpsermon/pkg"
)

var data *pkg.Data

func main() {

	// Setup PATH to include "/usr/local/bin" where brew installs things like youtube-dl and ffmpeg
	os.Setenv("PATH", fmt.Sprintf("%s:/usr/local/bin", os.Getenv("PATH")))

	// MAIN application window
	application := app.New()
	window := application.NewWindow("WPC Sermon")
	image := canvas.NewImageFromFile("assets/WPC_logo_brown_stacked.png")
	image.FillMode = canvas.ImageFillContain

	// Data
	data = pkg.NewData()
	defer data.Remove()

	// Setup image and input fields in a box
	input := container.New(layout.NewGridLayoutWithColumns(2),
		container.New(layout.NewMaxLayout(), image),
		container.New(layout.NewVBoxLayout(), InputForm(window)),
	)

	// Input input
	window.SetContent(container.NewVBox(
		input,
		container.NewVBox(
			widget.NewLabel("Downloading"),
			data.VideoProgress,
		),
		container.NewVBox(
			widget.NewLabel("Trimming"),
			data.TrimProgress,
		),
		container.NewVBox(
			widget.NewLabel("Finalizing"),
			data.FinalProgress,
		),
	))

	window.Resize(fyne.NewSize(1000, 500))
	window.CenterOnScreen()
	window.ShowAndRun()
	os.Exit(0)
}

func InputForm(window fyne.Window) *widget.Form {

	url := widget.NewEntry()
	start := widget.NewEntry()
	stop := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "URL",
				Widget: url,
			},
			{
				Text:   "Start Time",
				Widget: start,
			},
			{
				Text:   "Stop Time",
				Widget: stop,
			},
		},
	}

	form.OnSubmit = func() {

		// Disable the form now that Submit was pressed.
		form.Disable()

		data.URL = url.Text
		data.Start = start.Text
		data.Stop = stop.Text

		pkg.DownloadFile(data)
		pkg.Trim(data)
		pkg.Finalize(data)

		widget.ShowModalPopUp(
			widget.NewButton("\n\n        Finished       \n\n", func() { fyne.CurrentApp().Quit() }),
			window.Canvas(),
		)
	}

	return form
}
