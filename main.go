package main

import (
	"chess/engine"
	"chess/magic"
	"chess/uci"
)

func main() {
	err := magic.Prepare()
	if err != nil {
		panic(err)
	}
	engine.Run()
	uci.Start()
}
