package main

import (

	// STDLIB
	"os"

	// External
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"

	// Internal
	"github.com/marshall-mcmullen/wpsermon/pkg"
)

var data *pkg.Data

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	//_, filename, _, _ := runtime.Caller(0)

	// MAIN application window
	application := app.New()
	window := application.NewWindow("main")
	window.CenterOnScreen()
	image := canvas.NewImageFromFile("assets/WPC_logo_brown_stacked.png")
	image.FillMode = canvas.ImageFillOriginal

	// Data
	data = pkg.NewData()
	defer data.Remove()

	// Input prompt
	window.SetContent(container.NewVBox(
		image,
		InputForm(window),
		container.NewVBox(
			widget.NewLabel("Downloading Audio/Video"),
			data.AudioProgress,
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

	window.ShowAndRun()
	os.Exit(0)
}

func InputForm(window fyne.Window) *widget.Form {

	url := widget.NewEntry()
	start := widget.NewEntry()
	stop := widget.NewEntry()


	return &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: "URL",
				Widget: url,

			},
			{
				Text: "Start Time",
				Widget: start,
			},
			{
				Text: "Stop Time",
				Widget: stop,
			},
		},
		OnSubmit: func() {

			data.InputUrl = url.Text
			data.Start = start.Text
			data.Stop = stop.Text

			urls := pkg.GetURLs(data.InputUrl)
			data.VideoUrl = urls[0]
			data.AudioUrl = urls[1]
			log.WithFields(log.Fields{
				"data": data,
			}).Info("Data")

			pkg.DownloadAVFiles(data)
			pkg.Trim(data)
			pkg.Finalize(data)

			widget.ShowModalPopUp(
				widget.NewButton("Finished", func() { fyne.CurrentApp().Quit() }),
				window.Canvas(),
			)
		},
	}
}
