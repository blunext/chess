package main

import (
	"chess/generator"
	"chess/magic"
	"chess/uci"
)

func main() {
	err := magic.Prepare()
	if err != nil {
		panic(err)
	}
	generator.NewGenerator()
	uci.Start()
}
