package pkg

import (
	// STDLIB
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	// External
	"fyne.io/fyne/v2/widget"
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

func monitorProgress(cmd *exec.Cmd, data *Data, progress *widget.ProgressBar, totalFrames float64) {

	exp := regexp.MustCompile("^frame=([0-9]+)")

	pipe, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	CheckError(err)
	defer pipe.Close()

	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		fmt.Print(line)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
				break
			}

			CheckError(err)
		}

		matches := exp.FindStringSubmatch(line)

		if len(matches) == 2 {
			frameStr := matches[1]
			frame, err := strconv.ParseFloat(frameStr, 64)
			CheckError(err)

			percent := float64(frame) / float64(totalFrames)
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
		"-to", data.Stop,
		"-i", data.VideoPath,
		"-c", "copy",
		"-progress", "/dev/stdout",
		data.TrimmedPath,
	)
	printBanner(cmd.String())

	waitGroup.Add(1)
	go func() {
		err := cmd.Run()
		CheckError(err)
		waitGroup.Done()
	}()

	monitorProgress(cmd, data, data.TrimProgress, totalFrames)

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
		"-filter_complex", "[0:v]scale=1280x720[v0];[1:v]scale=1280x720[v1];[2:v]scale=1280x720[v2];[v0][0:a][v1][1:a][v2][2:a]concat=n=3:v=1:a=1[v][a]",
		"-map", "[v]",
		"-map", "[a]",
		"-progress", "/dev/stdout",
		"-f", "mp4",
		data.FinalVideoPath+".tmp",
	)
	printBanner(cmd.String())

	// Extract Audio only
	cmd2 := exec.Command("ffmpeg",
		"-i", data.FinalVideoPath,
		"-vn",
		"-acodec", "mp3",
		"-f", "mp3",
		data.FinalAudioPath+".tmp",
	)
	printBanner(cmd.String())

	waitGroup.Add(1)
	go func() {

		// Video
		os.Remove(data.FinalVideoPath + ".tmp")
		err := cmd.Run()
		CheckError(err)
		os.Rename(data.FinalVideoPath+".tmp", data.FinalVideoPath)

		// Audio
		os.Remove(data.FinalAudioPath + ".tmp")
		err = cmd2.Run()
		CheckError(err)
		os.Rename(data.FinalAudioPath+".tmp", data.FinalAudioPath)

		waitGroup.Done()
	}()

	monitorProgress(cmd, data, data.FinalProgress, totalFrames)

	waitGroup.Wait()
}
