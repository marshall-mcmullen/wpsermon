package pkg

import (
	// STDLIB
	"fmt"

	// External
	"github.com/fatih/color"
)

var printBoldGreen = color.New(color.FgGreen).Add(color.Bold).PrintfFunc()

func printBanner(text string) {
	printBoldGreen(fmt.Sprintf(" ● %v \n", text))
}
