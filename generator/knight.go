package generator

import "fmt"

func generateKnightMoves() squareMoves {
	var squareMoves = make(squareMoves)
	for pos := 0; pos < 64; pos++ {
		rank, file := rankAndFile(pos)

		fmt.Println(rank, file)

	}
	return squareMoves
}
