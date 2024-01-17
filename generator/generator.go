package generator

import (
	"chess/board"
)

type Moves []board.Bitboard
type PossibleMoves map[board.Bitboard]Moves

type MovesArray map[board.Bitboard][]Moves //slider list of moves

type Sliders map[board.Piece]MovesArray

type Generic map[board.Piece]PossibleMoves //King and Knight moves

func NewGenerator() (Sliders, Generic) {
	sliders := make(Sliders)
	sliders[board.Rook] = generateRookMoves()
	sliders[board.Bishop] = generateBishopMoves()
	sliders[board.Queen] = generateQueenMoves(sliders[board.Rook], sliders[board.Bishop])

	generic := make(Generic)
	generic[board.King] = kingMoves()
	generic[board.Knight] = knightMoves()

	return sliders, generic
}

func generateQueenMoves(rookMoves, bishopMoves MovesArray) MovesArray {
	var squareMoves = make(MovesArray)
	var pos board.Bitboard
	for pos = 0; pos < 64; pos++ {
		var directions []Moves
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
