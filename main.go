package main

import (
	"os"

	"github.com/adil-chbada/extract-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}