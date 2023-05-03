package main

import (
	"chess/generator"
	"chess/magic"
)

func main() {
	err := magic.Prepare()
	if err != nil {
		panic(err)
	}
	generator.NewGenerator()
}
