package generator

import (
	"chess/board"
)

func generateRookMoves() board.SquareMoves {
	var squareMoves = make(board.SquareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions [][]board.Bitboard
		moves := rookDown(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookUp(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookRight(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookLeft(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}

		squareMoves[board.IndexToBitBoard(pos)] = exactSize(directions)
	}
	return squareMoves
}

func rookDown(pos int) []board.Bitboard {
	var list []board.Bitboard
	for i := 1; i <= rank(pos); i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookUp(pos int) []board.Bitboard {
	var list []board.Bitboard
	for i := 1; i <= 7-rank(pos); i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookRight(pos int) []board.Bitboard {
	var list []board.Bitboard
	file := file(pos) + 1
	for i := 1; i <= 8-file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookLeft(pos int) []board.Bitboard {
	var list []board.Bitboard
	file := file(pos) + 1
	for i := 1; i < file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i)
		list = append(list, newPos)
	}
	return exactSize(list)
}

// GenerateRookMovesForTesting exports the rook move generator for testing purposes.
func GenerateRookMovesForTesting() board.SquareMoves {
	return generateRookMoves()
}
