package pkg

import (
	// STDLIB
	"io/ioutil"
	"os"
	"path/filepath"

	// External

	"fyne.io/fyne/v2/widget"
)

type Data struct {
	TempDir string
	URL     string
	Start   string
	Stop    string

	// Paths
	VideoPath        string
	AudioPath        string
	TrimProgressPath string
	TrimmedPath      string
	IntroPath        string
	EndingPath       string
	FinalVideoPath   string
	FinalAudioPath   string

	// GUI Elements
	VideoProgress *widget.ProgressBar
	TrimProgress  *widget.ProgressBar
	FinalProgress *widget.ProgressBar
}

func NewData() *Data {

	dir, err := ioutil.TempDir("", "wpsermon")
	CheckError(err)

	return &Data{
		TempDir:          dir,
		IntroPath:        filepath.Join("assets", "intro.mp4"),
		EndingPath:       filepath.Join("assets", "ending.mp4"),
		FinalVideoPath:   filepath.Join(os.Getenv("HOME"), "Desktop", "sermon.mp4"),
		FinalAudioPath:   filepath.Join(os.Getenv("HOME"), "Desktop", "sermon.mp3"),
		AudioPath:        filepath.Join(dir, "audio.mp4"),
		VideoPath:        filepath.Join(dir, "video.mp4"),
		VideoProgress:    widget.NewProgressBar(),
		TrimProgress:     widget.NewProgressBar(),
		TrimmedPath:      filepath.Join(dir, "trimmed.mp4"),
		TrimProgressPath: filepath.Join(dir, "trim_progress.txt"),
		FinalProgress:    widget.NewProgressBar(),
	}
}

func (data *Data) Remove() {
	os.RemoveAll(data.TempDir)
}
