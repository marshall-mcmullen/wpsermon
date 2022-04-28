package main

import (

	// STDLIB
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	// External
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	log "github.com/sirupsen/logrus"

	// Internal
	"github.com/marshall-mcmullen/wpsermon/pkg"
)

type data struct {
	inputUrl string
	start string
	stop string
	videoUrl string
	audioUrl string
}

func checkError(err error) {

	if err != nil {

		gui := app.New()
		window := gui.NewWindow("Error")
		content := widget.NewLabel(err.Error())

		window.SetContent(content)
		window.ShowAndRun()
		os.Exit(1)
	}
}

func main() {

	// Setup
	log.SetFormatter(&log.JSONFormatter{})

	application := app.New()
	window := application.NewWindow("main")

	url := widget.NewEntry()
	start := widget.NewEntry()
	stop := widget.NewEntry()
	data := &data{}

	form := &widget.Form{
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

			data.inputUrl = url.Text
			data.start =  start.Text
			data.stop = stop.Text

			window.Hide()

			urls := getURLs(data.inputUrl)
			data.videoUrl = urls[0]
			data.audioUrl = urls[1]

			log.WithFields(log.Fields{
				"data": data,
			}).Info("Data")

			process(data)

			window.Close()
		},
	}

	window.Resize(fyne.NewSize(600, 150))
	window.SetContent(form)
	window.ShowAndRun()
}

func process(data *data) {

	// Download Video and Audio URLs
	dir, err := ioutil.TempDir("", "wpsermon")
	checkError(err)
	defer os.RemoveAll(dir)

	pkg.DownloadFile(filepath.Join(dir, "video.mp4"), data.audioUrl)
}

func prompt() *data {

	window := fyne.CurrentApp().NewWindow("Whispering Pines Church")

	url := widget.NewEntry()
	start := widget.NewEntry()
	stop := widget.NewEntry()
	data := &data{}

	form := &widget.Form{
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

			data.inputUrl = url.Text
			data.start =  start.Text
			data.stop = stop.Text

			process(data)
			window.Close()
		},
	}

	window.Resize(fyne.NewSize(600, 150))
	window.SetContent(form)
	window.Show()

	return data
}

func getURLs(url string) []string{

	cmd := exec.Command("youtube-dl", "-g", url)
	stdout, err := cmd.Output()
	checkError(err)

	log.Info(string(stdout))

	return strings.Split(string(stdout), "\n")
}

func progress() {
	myApp := app.New()
	window := myApp.NewWindow("ProgressBar Widget")

	progress := widget.NewProgressBar()
	infinite := widget.NewProgressBarInfinite()

	go func() {
		for i := 0.0; i <= 1.0; i += 0.1 {
			time.Sleep(time.Millisecond * 250)
			progress.SetValue(i)
		}
	}()

	window.SetContent(container.NewVBox(progress, infinite))
	window.ShowAndRun()
}
