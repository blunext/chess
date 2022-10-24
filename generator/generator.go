package generator

import (
	"chess/board"
)

type possibleMoves []uint64

type sliderMoves map[board.Bitboard][]possibleMoves

type generatedMoves map[uint8]sliderMoves

type knightMoves map[uint64]possibleMoves

func NewGenerator() generatedMoves {
	moves := make(generatedMoves)
	moves[board.Rook] = generateRookMoves()
	moves[board.Bishop] = generateBishopMoves()
	moves[board.Queen] = generateQueenMoves(moves[board.Rook], moves[board.Bishop])
	return moves
}

func generateQueenMoves(rookMoves, bishopMoves sliderMoves) sliderMoves {
	var squareMoves = make(sliderMoves)
	var pos board.Bitboard
	for pos = 0; pos < 64; pos++ {
		var directions []possibleMoves
		for _, list := range rookMoves[pos] {
			directions = append(directions, list)
		}
		for _, list := range bishopMoves[pos] {
			directions = append(directions, list)
		}
		squareMoves[pos] = exactSize(directions)
	}
	return squareMoves
}

func exactSize[T any](list []T) []T {
	exactSlice := make([]T, len(list))
	copy(exactSlice, list)
	return exactSlice
}

func file(pos int) int {
	return pos & 7
}

func rank(pos int) int {
	return 7 - pos>>3
}

func rankAndFile(pos int) (int, int) {
	return rank(pos), file(pos)
}
