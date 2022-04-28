package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
	Current uint64
	Window fyne.Window
	Progress *widget.ProgressBar
	Infinite *widget.ProgressBarInfinite
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	percent := float64(wc.Current) / float64(wc.Total)
	fmt.Printf("\rDownloading... Current=%d Total=%d (%f %%)", wc.Current, wc.Total, percent)
	wc.Progress.SetValue(percent)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Head(url)
	if err != nil {
		return err
	}

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// Get the data
	resp, err = http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		Total: uint64(size),
		Window: fyne.CurrentApp().NewWindow("ProgressBar Widget"),
		Progress: widget.NewProgressBar(),
		Infinite: widget.NewProgressBarInfinite(),
	}

	counter.Window.SetContent(container.NewVBox(counter.Progress, counter.Infinite))
	counter.Window.Show()

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}

	counter.Window.Close()

	return nil
}

