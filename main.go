package main

import (
	"chess/generator"
	"chess/magic"
)

func main() {
	magic.Prepare()
	generator.NewGenerator()
}
