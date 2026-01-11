package generator

import (
	"chess/board"
)

func generateBishopMoves() board.SquareMoves {
	var squareMoves = make(board.SquareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions [][]board.Bitboard
		moves := bishopSE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopSW(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopNW(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopNE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}

		squareMoves[board.IndexToBitBoard(pos)] = exactSize(directions)
	}
	return squareMoves
}

func bishopSE(pos int) []board.Bitboard {
	var list []board.Bitboard
	rank, file := rankAndFile(pos)
	for {
		rank--
		file++
		if rank < 0 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit(rank*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopSW(pos int) []board.Bitboard {
	var list []board.Bitboard
	rank, file := rankAndFile(pos)
	for {
		rank--
		file--
		if rank < 0 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit(rank*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNE(pos int) []board.Bitboard {
	var list []board.Bitboard
	rank, file := rankAndFile(pos)
	for {
		rank++
		file++
		if rank >= 8 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit(rank*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNW(pos int) []board.Bitboard {
	var list []board.Bitboard
	rank, file := rankAndFile(pos)
	for {
		rank++
		file--
		if rank >= 8 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit(rank*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

// GenerateBishopMovesForTesting exports the bishop move generator for testing purposes.
func GenerateBishopMovesForTesting() board.SquareMoves {
	return generateBishopMoves()
}
