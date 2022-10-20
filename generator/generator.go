package generator

import (
	"fmt"

	"chess/board"
)

type possibleMoves []board.Bitboard

type squareMoves map[board.Bitboard][]possibleMoves

type generatedMoves map[uint8]squareMoves

func NewGenerator() generatedMoves {
	moves := make(generatedMoves)
	moves[board.Rook] = generateRookMoves()
	moves[board.Bishop] = generateBishopMoves()
	fmt.Println(len(moves))
	return moves
}
