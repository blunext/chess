package engine

import (
	"chess/board"
	"chess/generator"
	"fmt"
)

func Run() {
	position := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1")
	pieceMoves := generator.NewGenerator()

	// Generate moves using new Move struct
	moves := position.GenerateMoves(pieceMoves)

	fmt.Println("Generated sliding piece moves:")
	for _, m := range moves {
		fmt.Println(" ", m)
	}
	fmt.Printf("Total: %d moves\n", len(moves))
}
