package pkg

import (
	// STDLIB
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	// External
	"fyne.io/fyne/v2/widget"
)

// ProgressWriter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type ProgressWriter struct {
	Total    int
	Current  int
	Progress *widget.ProgressBar
}

func (wc *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += n
	wc.PrintProgress()
	return n, nil
}

func (wc ProgressWriter) PrintProgress() {
	percent := float64(wc.Current) / float64(wc.Total)
	wc.Progress.SetValue(percent)
}

type DownloadFile struct {
	URL      string
	Output   string
	Writer   *ProgressWriter
	Progress *widget.ProgressBar
}

func downloadSingleFile(file *DownloadFile) {

	out, err := os.Create(file.Output + ".tmp")
	CheckError(err)

	resp, err := http.Head(file.URL)
	CheckError(err)

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	file.Writer = &ProgressWriter{
		Total:    size,
		Progress: file.Progress,
	}

	// Get the data
	resp, err = http.Get(file.URL)
	CheckError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, io.TeeReader(resp.Body, file.Writer))
	CheckError(err)

	// Close the file without defer so it can happen before Rename()
	out.Close()
	err = os.Rename(file.Output+".tmp", file.Output)
	CheckError(err)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadAVFiles(data *Data) {

	var waitGroup sync.WaitGroup

	audio := &DownloadFile{
		URL:      data.AudioUrl,
		Output:   data.AudioPath,
		Progress: data.AudioProgress,
	}

	video := &DownloadFile{
		URL:      data.VideoUrl,
		Output:   data.VideoPath,
		Progress: data.VideoProgress,
	}

	waitGroup.Add(1)
	go func() {
		downloadSingleFile(audio)
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		downloadSingleFile(video)
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

func GetURLs(url string) []string {

	cmd := exec.Command("yt-dlp", "-g", url)
	stdout, err := cmd.Output()
	CheckError(err)

	return strings.Split(string(stdout), "\n")
}
