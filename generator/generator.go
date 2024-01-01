package generator

import (
	"chess/board"
)

type PossibleMoves map[board.Bitboard][]board.Bitboard

type GeneratedMoves map[uint8]PossibleMoves

func NewGenerator() GeneratedMoves {
	moves := make(GeneratedMoves)
	moves[board.Rook] = generateRookMoves()
	moves[board.Bishop] = generateBishopMoves()
	moves[board.Queen] = generateQueenMoves(moves[board.Rook], moves[board.Bishop])
	moves[board.Knight] = generateKnightMoves()
	return moves
}

func generateQueenMoves(rookMoves, bishopMoves PossibleMoves) PossibleMoves {
	var squareMoves = make(PossibleMoves)
	var pos board.Bitboard
	for pos = 0; pos < 64; pos++ {
		var directions []board.Bitboard
		directions = append(directions, rookMoves[pos]...)
		directions = append(directions, bishopMoves[pos]...)
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
