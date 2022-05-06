package pkg

import (
	// STDLIB
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	// External
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

func getTotalFrames(video string) float64 {

	// First we need to figure out the total number of frames in the video
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-count_packets",
		"-show_entries",
		"stream=nb_read_packets",
		"-of", "csv=p=0",
		video)

	var stdoutBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	err := cmd.Run()

	totalFrames, err := strconv.ParseFloat(strings.Replace(stdoutBuffer.String(), "\n", "", -1), 64)
	CheckError(err)

	return totalFrames
}

func monitorProgress(cmd *exec.Cmd, progress *widget.ProgressBar, totalFrames float64) {

	exp := regexp.MustCompile("^frame=([0-9]+)")

	pipe, _ := cmd.StdoutPipe()
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		line, err = reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		CheckError(err)

		matches := exp.FindStringSubmatch(line)

		if len(matches) == 2 {
			frameStr := matches[1]
			frame, err := strconv.ParseFloat(frameStr, 64)
			CheckError(err)

			percent := float64(frame) / float64(totalFrames)
			log.Infof("Frame: %f/%f Percent=%f", frame, totalFrames, percent)
			progress.SetValue(percent)
		}
	}

	progress.SetValue(100)
}

func Trim(data *Data) {

	totalFrames := getTotalFrames(data.VideoPath)

	var waitGroup sync.WaitGroup

	cmd := exec.Command("ffmpeg",
		"-ss", data.Start,
		"-i", data.AudioPath,
		"-to", data.Stop,
		"-ss", data.Start,
		"-i", data.VideoPath,
		"-to", data.Stop,
		"-c", "copy",
		"-progress", "/dev/stdout",
		data.TrimmedPath,
	)

	waitGroup.Add(1)
	go func() {
		err := cmd.Run()
		CheckError(err)
		waitGroup.Done()
	}()

	monitorProgress(cmd, data.TrimProgress, totalFrames)

	waitGroup.Wait()
}

func Finalize(data *Data) {

	totalFrames := getTotalFrames(data.IntroPath)
	totalFrames += getTotalFrames(data.TrimmedPath)
	totalFrames += getTotalFrames(data.EndingPath)

	var waitGroup sync.WaitGroup

	// Remaster Audio and Video and concatenate all streams
	cmd := exec.Command("ffmpeg",
		"-vsync", "0",
		"-i", data.IntroPath,
		"-i", data.TrimmedPath,
		"-i", data.EndingPath,
		"-filter_complex", "[0:v:0][0:a:0][1:v:0][1:a:0][2:v:0][2:a:0]concat=n=3:v=1:a=1[outv][outa]",
		"-map", "[outv]",
		"-map", "[outa]",
		"-progress", "/dev/stdout",
		"-f", "mp4",
		data.FinalVideoPath+".tmp",
	)

	// Extract Audio only
	cmd2 := exec.Command("ffmpeg",
		"-i", data.FinalVideoPath,
		"-vn",
		"-acodec", "mp3",
		"-f", "mp3",
		data.FinalAudioPath+".tmp")

	waitGroup.Add(1)
	go func() {

		// Video
		err := cmd.Run()
		CheckError(err)
		os.Rename(data.FinalVideoPath+".tmp", data.FinalVideoPath)

		// Audio
		err = cmd2.Run()
		CheckError(err)
		os.Rename(data.FinalAudioPath+".tmp", data.FinalAudioPath)

		waitGroup.Done()
	}()

	monitorProgress(cmd, data.FinalProgress, totalFrames)

	waitGroup.Wait()
}
