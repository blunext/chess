package board

import (
	"fmt"
	"testing"
)

func TestPosition_Pretty(t *testing.T) {
	position := CreatePositionFormFEN(InitialPosition)
	fmt.Println("\nInitial position:")
	fmt.Println(position.Pretty())

	position2 := CreatePositionFormFEN("rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1")
	fmt.Println("\nAfter 1.d4 e6:")
	fmt.Println(position2.Pretty())
}
