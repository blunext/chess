package generator

import (
	"chess/board"
)

type SquareMoves map[board.Bitboard][]board.Bitboard
type Generics map[board.Piece]SquareMoves
type SliderSquareMoves map[board.Bitboard][][]board.Bitboard
type Sliders map[board.Piece]SliderSquareMoves

func NewGenerator() (Sliders, Generics) {
	sliders := make(Sliders)
	sliders[board.Rook] = generateRookMoves()
	sliders[board.Bishop] = generateBishopMoves()
	sliders[board.Queen] = generateQueenMoves(sliders[board.Rook], sliders[board.Bishop])

	generic := make(Generics)
	generic[board.King] = kingMoves()
	generic[board.Knight] = knightMoves()

	return sliders, generic
}

func generateQueenMoves(rookMoves, bishopMoves SliderSquareMoves) SliderSquareMoves {
	var squareMoves = make(SliderSquareMoves)
	var pos board.Bitboard
	for pos = 0; pos < 64; pos++ {
		var directions [][]board.Bitboard
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
