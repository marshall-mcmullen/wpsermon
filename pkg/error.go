package pkg

import (
	// STDLIB
	"os"

	// External
	log "github.com/sirupsen/logrus"
)

func CheckError(err error) {

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

