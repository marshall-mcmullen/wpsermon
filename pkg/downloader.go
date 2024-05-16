package pkg

import (
	// STDLIB
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func DownloadFile(data *Data) {

	cmd := exec.Command("yt-dlp",
		"--newline",
		"-f", "bestvideo+bestaudio",
		"--merge-output-format", "mp4",
		"-o", data.VideoPath,
		data.URL,
	)
	printBanner(cmd.String())

	// Setup pipe
	stdout, err := cmd.StdoutPipe()
	CheckError(err)

	// Start program
	err = cmd.Start()
	CheckError(err)

	// Create a string scanner to read the output from the pipe
	scanner := bufio.NewScanner(stdout)

	// Create a regex pattern to match progress lines
	// [download]  10.0% of    1.47GiB at   92.78MiB/s ETA 00:14
	progressPattern := regexp.MustCompile(`\[download\]\s+(\d+\.\d+)%\s+of\s+(\d+\.\d+\w+)\s+at\s+(\d+\.\d+\w+/s)\s+ETA\s+(\d+:\d+)`)

	// Create a goroutine to read from the pipe and update the progress label
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if strings.Contains(line, "[download]") {
			match := progressPattern.FindStringSubmatch(line)
			if len(match) == 5 {
				percentStr := match[1]
				percent, err := strconv.ParseFloat(percentStr, 64)
				CheckError(err)
				data.VideoProgress.SetValue(percent / 100.0)
			}
		}
	}

	CheckError(scanner.Err())

	err = cmd.Wait()
	CheckError(err)
}
