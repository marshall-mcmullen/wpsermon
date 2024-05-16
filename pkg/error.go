package pkg

import (
	"fmt"
	"os"
)

func CheckError(err error) {

	if err != nil {
		panic(err)
	}
}

func CheckErrorWithOutput(err error, output string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Fatal Error: %v: %v", err, output))
		panic(err)
	}
}
