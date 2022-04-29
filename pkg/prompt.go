package pkg

import (

	// STDLIB
	"os/exec"
	"strings"

	// External
	log "github.com/sirupsen/logrus"
)

func GetURLs(url string) []string{

	cmd := exec.Command("youtube-dl", "-g", url)
	stdout, err := cmd.Output()
	CheckError(err)

	log.Info(string(stdout))

	return strings.Split(string(stdout), "\n")
}
