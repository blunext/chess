package generator

import (
	"chess/board"
)

func NewGenerator() board.PieceMoves {
	pieceMoves := make(board.PieceMoves)
	pieceMoves[board.Rook] = generateRookMoves()
	pieceMoves[board.Bishop] = generateBishopMoves()
	pieceMoves[board.Queen] = generateQueenMoves(pieceMoves[board.Rook], pieceMoves[board.Bishop])
	pieceMoves[board.King] = kingMoves()
	pieceMoves[board.Knight] = knightMoves()

	return pieceMoves
}

func generateQueenMoves(rookMoves, bishopMoves board.SquareMoves) board.SquareMoves {
	var squareMoves = make(board.SquareMoves)
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
