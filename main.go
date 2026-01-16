package main

import (
	"os"

	"chess/engine"
	"chess/magic"
	"chess/uci"
)

func main() {
	err := magic.Prepare()
	if err != nil {
		panic(err)
	}

	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "play":
			engine.Play()
			return
		case "uci":
			uci.Start()
			return
		}
	}

	// Default: UCI mode
	uci.Start()
}
