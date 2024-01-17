package engine

import (
	"chess/board"
	"chess/generator"
)

func Run() {
	position := board.CreatePositionFormFEN(board.InitialPosition)
	sliders, _ := generator.NewGenerator()
	position.AllBishops(sliders)
}
