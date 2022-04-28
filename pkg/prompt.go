package pkg

import (

	// STDLIB
	"os/exec"
	"strings"
	"sync"

	// External
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	InputUrl string
	Start string
	Stop string
	VideoUrl string
	AudioUrl string
	VideoPath string
	AudioPath string
}

func Prompt() *Data {

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	window := fyne.CurrentApp().NewWindow("Whispering Pines Church")
	window.CenterOnScreen()

	url := widget.NewEntry()
	start := widget.NewEntry()
	stop := widget.NewEntry()
	data := &Data{}

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

			data.InputUrl = url.Text
			data.Start =  start.Text
			data.Stop = stop.Text
			window.Close()
			waitGroup.Done()
		},
	}

	window.Resize(fyne.NewSize(600, 150))
	window.SetContent(form)
	window.Show()

	// Wait for GUI prompting to complete
	waitGroup.Wait()

	urls := getURLs(data.InputUrl)
	data.VideoUrl = urls[0]
	data.AudioUrl = urls[1]

	log.WithFields(log.Fields{
		"data": data,
	}).Info("Data")

	return data
}

func getURLs(url string) []string{

	cmd := exec.Command("youtube-dl", "-g", url)
	stdout, err := cmd.Output()
	CheckError(err)

	log.Info(string(stdout))

	return strings.Split(string(stdout), "\n")
}
