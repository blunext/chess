package engine

import (
	"chess/board"
	"chess/generator"
)

func Run() {
	//position := board.CreatePositionFormFEN(board.InitialPosition)
	position := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1")
	sliders, _ := generator.NewGenerator()
	//position.AllBishops(sliders, board.Bishop)
	position.AllBishops(sliders)
}
