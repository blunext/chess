package generator

import (
	"chess/board"
)

func NewGenerator() (board.Sliders, board.Generics) {
	sliders := make(board.Sliders)
	sliders[board.Rook] = generateRookMoves()
	sliders[board.Bishop] = generateBishopMoves()
	sliders[board.Queen] = generateQueenMoves(sliders[board.Rook], sliders[board.Bishop])

	generic := make(board.Generics)
	generic[board.King] = kingMoves()
	generic[board.Knight] = knightMoves()

	return sliders, generic
}

func generateQueenMoves(rookMoves, bishopMoves board.SliderSquareMoves) board.SliderSquareMoves {
	var squareMoves = make(board.SliderSquareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions [][]board.Bitboard
		bbIndex := board.IndexToBitBoard(pos)
		directions = append(directions, rookMoves[bbIndex]...)
		directions = append(directions, bishopMoves[bbIndex]...)
		squareMoves[bbIndex] = exactSize(directions)
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
