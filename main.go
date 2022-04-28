package main

import (

	// STDLIB
	"io/ioutil"
	"os"
	"path/filepath"

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

func main() {

	// Setup
	log.SetFormatter(&log.JSONFormatter{})

	application := app.New()
	window := application.NewWindow("main")
	window.CenterOnScreen()

	image := canvas.NewImageFromFile("assets/WPC_logo_green_stacked.png")
	image.FillMode = canvas.ImageFillOriginal
	button := widget.NewButton("Download Sermon", func() { process(window); })

	window.SetContent(container.NewVBox(image, button))
	window.ShowAndRun()
	os.Exit(0)
}

func process(window fyne.Window) {

	dir, err := ioutil.TempDir("", "wpsermon")
	pkg.CheckError(err)
	defer os.RemoveAll(dir)

	// PROMPT
	window.Hide()
	data := pkg.Prompt()
	data.AudioPath = filepath.Join(dir, "audio.mp4")
	data.VideoPath = filepath.Join(dir, "video.mp4")

	// DOWNLOAD
	window.Hide()
	pkg.DownloadAVFiles(data)

	// MODIFY
	window.Hide()
	pkg.Trim(data)
}
